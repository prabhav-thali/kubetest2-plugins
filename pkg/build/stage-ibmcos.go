package build

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"

	"k8s.io/klog/v2"
)

type IBMCOSStager struct {
	StageLocation   string
	RepoRoot        string
	Region          string
	Bucket          string
	Path            string
	TargetBuildArch string
	Credentials     *credentials.Credentials
}

func NewIBMCOSStager(stageLocation, repoRoot, targetBuildArch, cosCredType string) (*IBMCOSStager, error) {
	re := regexp.MustCompile(`^([a-zA-Z]+):\/\/([a-zA-Z0-9-]+)\/([a-zA-Z0-9-]+)(\/.*)?$`)
	matches := re.FindStringSubmatch(stageLocation)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid IBM COS stagelocation, missing region, bucket information, expected format is cos://us/bucket123/<PATH>")
	}

	var cred *credentials.Credentials

	switch cosCredType {
	case "shared":
		cred = credentials.NewSharedCredentials("", "")
	case "cos_hmac":
		cred = NewCosHmacCredentials("")
	default:
		return nil, errors.New("invalid credential type: " + cosCredType)
	}

	return &IBMCOSStager{
		StageLocation:   stageLocation,
		RepoRoot:        repoRoot,
		Region:          matches[2],
		Bucket:          matches[3],
		Path:            matches[4],
		TargetBuildArch: targetBuildArch,
		Credentials:     cred,
	}, nil
}

var _ Stager = &IBMCOSStager{}

func (i *IBMCOSStager) getS3Client() *s3.S3 {
	conf := aws.NewConfig().
		WithRegion(fmt.Sprintf("%s-standard", i.Region)).
		WithEndpoint(fmt.Sprintf("https://s3.%s.cloud-object-storage.appdomain.cloud", i.Region)).
		WithCredentials(i.Credentials).
		WithS3ForcePathStyle(true)

	sess := session.Must(session.NewSession())
	return s3.New(sess, conf)
}

// Stage implements Stager.
func (i *IBMCOSStager) Stage(version string) error {
	client := i.getS3Client()
	tgzFile := "kubernetes-server-" + strings.ReplaceAll(i.TargetBuildArch, "/", "-") + ".tar.gz"
	destinationKey := aws.String(i.Path + "/" + version + "/" + tgzFile)
	klog.Infof("uploading %s to location %s/%s", tgzFile, i.StageLocation, version)

	f, err := os.Open(i.RepoRoot + "/_output/release-tars/" + tgzFile)
	if err != nil {
		return err
	}
	defer f.Close()

	uploader := s3manager.NewUploaderWithClient(client)

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(i.Bucket),
		Key:    destinationKey,
		Body:   aws.ReadSeekCloser(bufio.NewReader(f)),
	}
	if _, err = uploader.Upload(upParams); err != nil {
		return err
	}
	klog.Info("file uploaded")
	return nil
}

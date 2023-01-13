package baidu

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *BaiduCloud) ListObjects(endpoint, bucketName, prefix string) (objects []cloud.ObjectProperties, err error) {
	if res, err := p.bosClient.ListBuckets(); err != nil {
		return nil, err
	} else {
		for _, buck := range res.Buckets {
			if buck.Location == endpoint && buck.Name == bucketName {
				listObjectResult, err := p.bosClient.ListObjects(bucketName, nil)
				if err != nil {
					return nil, err
				}
				for _, obj := range listObjectResult.Contents {
					object := cloud.ObjectProperties{
						Name: obj.Key,
					}
					objects = append(objects, object)
				}
				break
			}
		}
	}
	return
}
func (p *BaiduCloud) ListBucket(endpoint string) ([]cloud.BucketProperties, error) {
	buckets := []cloud.BucketProperties{}
	if res, err := p.bosClient.ListBuckets(); err != nil {
		return nil, err
	} else {
		for _, b := range res.Buckets {
			if b.Location == endpoint {
				bucket := cloud.BucketProperties{
					Name: b.Name,
				}
				buckets = append(buckets, bucket)
			}
		}
	}
	return buckets, nil
}

func (p *BaiduCloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {
	// todo
	return ""
}

func (p *BaiduCloud) GetObjectDownloadUrl(bucketName, objectName string) (string, error) {
	url := p.bosClient.BasicGeneratePresignedUrl(bucketName, objectName, 300)
	return url, nil
}

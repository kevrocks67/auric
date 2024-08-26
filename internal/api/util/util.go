package util

import (
	"strings"
)

/* ExtractNameFromGoldenUri returns the artifact name
 *
 * Given a URI such as /golden/qcow2/rocky9-base/prod
 * the returned value would be rocky9-base
 */
func ExtractNameFromGoldenUri(uri string) string {
	return strings.Split(uri, "/")[3]
}

/* ExtractTypeFromGoldenUri returns the artifact name
 *
 * Given a URI such as /golden/qcow2/rocky9-base/prod
 * the returned value would be qcow2
 */
func ExtractTypeFromGoldenUri(uri string) string {
	return strings.Split(uri, "/")[2]
}

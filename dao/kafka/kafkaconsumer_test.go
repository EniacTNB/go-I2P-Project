package kafka_test

import (
	"gin-vue/dao/kafka"
	"testing"
)

func TestParseRouterInfo(t *testing.T) {

	// testKafkaMesg := "{\"pubkey\": \"BtLD~gAH1vP2O1I1cZ87MHhrMRVYaPmxB0owlflXTP8=\", \"signkey\": \"~gkz-04bNOWHKT56EihMR0umin2q2z7Vf7WwPgLu9GU=\", \"options\": {\"caps\": \"NfR\", \"netId\": \"2\", \"netdb.knownLeaseSets\": \"37\", \"netdb.knownRouters\": \"3859\", \"router.version\": \"0.9.31\"}, \"addrs\": [{\"cost\": 4, \"transport\": \"SSU\", \"options\": {\"caps\": \"BC\", \"host\": \"114.214.191.255\", \"key\": \"3WxV3BXFYF2vxI~CXKC4WNkEVomB1tntZ5p2lC6kKOo=\", \"mtu\": \"1488\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}, {\"cost\": 6, \"transport\": \"SSU\", \"options\": {\"caps\": \"BC\", \"host\": \"2.219.3.198\", \"key\": \"3WxV3BXFYF2vxI~CXKC4WNkEVomB1tntZ5p2lC6kKOo=\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}, {\"cost\": 9, \"transport\": \"NTCP\", \"options\": {\"host\": \"2a02:c7d:202f:ae00:bdb2:fb70:8aeb:90a\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}, {\"cost\": 10, \"transport\": \"NTCP\", \"options\": {\"host\": \"2.219.3.198\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}], \"cert\": {\"signature_type\": \"EdDSA_SHA512_Ed25519\", \"crypto_type\": \"ElGamal\"}, \"published\": 1635859281314, \"signature\": \"Pj~JQgamEKDGHHgDqTgaUZxpdZXrrpVeTcFFstJ1CfE=\",\"filename\":\"routerInfo-JaMocLiyYtoeP~jw7EllxiVIM-o3xivw93PLtgkn2Go=.dat\"}"
	// kafka.ParseRouterInfo(testKafkaMesg)
	kafka.RouterInfoLinstenr()

}

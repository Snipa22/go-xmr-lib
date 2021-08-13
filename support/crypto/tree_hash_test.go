package crypto

import (
	"encoding/hex"
	"testing"
)

func TestTreeHash(t *testing.T) {
	var testVals [][32]byte
	var tempBuf [32]byte

	hexBuf, _ := hex.DecodeString("676567f8b1b470207c20d8efbaacfa64b2753301b46139562111636f36304bb8")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)
	result := TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "676567f8b1b470207c20d8efbaacfa64b2753301b46139562111636f36304bb8" {
		t.Fatal("Unable to hash a single length object properly")
	}

	hexBuf, _ = hex.DecodeString("3124758667bc8e76e25403eee75a1044175d58fcd3b984e0745d0ab18f473984")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[0] = tempBuf
	hexBuf, _ = hex.DecodeString("975ce54240407d80eedba2b395bcad5be99b5c920abc2423865e3066edd4847a")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)
	result = TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "5077570fed2363a14fa978218185b914059e23517faf366f08a87cf3c47fd58e" {
		t.Fatalf("Unable to hash a double length object properly, received %v", hex.EncodeToString(result[:]))
	}

	hexBuf, _ = hex.DecodeString("decc1e0aa505d7d5fbe8ed823d7f5da55307c4cc7008e306da82dbce492a0576")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[0] = tempBuf
	hexBuf, _ = hex.DecodeString("dbcf0c26646d36b36a92408941f5f2539f7715bcb1e2b1309cedb86ae4211554")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[1] = tempBuf
	hexBuf, _ = hex.DecodeString("f56f5e6b2fce16536e44c851d473d1f994793873996ba448dd59b3b4b922b183")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)

	result = TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "f8e26aaa7c36523cea4c5202f2df159c62bf70d10670c96aed516dbfd5cb5227" {
		t.Fatalf("Unable to hash a triple length object properly, received %v", hex.EncodeToString(result[:]))
	}

	hexBuf, _ = hex.DecodeString("53edbbf98d3fa50a85fd2d46c42502aafad3fea30bc25ba4f16ec8bf4a475c4d")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[0] = tempBuf
	hexBuf, _ = hex.DecodeString("87da8ad3e5c90aae0b10a559a77a0985608eaa3cc3dd338239be52572c3bdf4b")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[1] = tempBuf
	hexBuf, _ = hex.DecodeString("a403d27466991997b3cf4e8d238d002a1451ccc9c4790269d0f0085d9382d60f")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[2] = tempBuf
	hexBuf, _ = hex.DecodeString("ef37717f59726e4cc8787d5d2d75238ba9adb9627a8f4aeeec8d80465ed3f5fb")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)

	result = TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "45f6e06fc0263e667caddd8fba84c9fb723a961a01a5b115f7cab7fe8f2c7e44" {
		t.Fatalf("Unable to hash a quad length object properly, received %v", hex.EncodeToString(result[:]))
	}

	testVals = [][32]byte{}
	hexBuf, _ = hex.DecodeString("21f750d5d938dd4ed1fa4daa4d260beb5b73509de9a9b145624d3f1afb671461b07d768cf1f5f8266b89ecdc150a2ad55ccd76d4c12d3a380b21862809a85af623269a23ee1b4694b26aa317b5cd4f259925f6b3288a8f60fb871b1ad3ac00cb1e6c55eddfc438e1f3e7b638ea6026cc01495010bafdfd789c47dff282c1af4c6a8f83e5f2fca6940a756ef4faa15c7137082a7c31dffe0b2f5112d126ad4af1d536c0e626cc9d2fe1b72256f5285728558f22a3dbb36e0918bcfc01d4ae7284d0bfb8e90647cdb01c292a53a31ff3fe6f350882f1dae2b09374db45f4d54c67d3b4e0829c4f9f63ad235d8ef838d8fb39546d90d99bbd831aff55dbbb642e2bf529ceccd0479b9f194475c2a15143f0edac762e9bbce810436e765550c69e234c22276c41d7d7e28c10afc5e144a9ce32aa9c0f28bb4fcf171af7d7404fa5e28b79dc97bd4147f4df6d38b935bd83fb634414bae9d64a32ab45384fba5b8da5c147d51cd2a8f7f2a9c07b1bddc5b28b74bf0c0f0632ac2fc43d0d306dd1ac1481cabe60a358d6043d4733202d489664a929d6bf76a39828954846beb47a3baacb35d2065cbe3ad34cf78bf895f6323a6d76fc1256306f58e4baecabd7a779388c6bf2734897c193d39c343fce49a456f0ef84cf963593c5401a14621cc6ec1bef01b53735ccb02bc96c5fd454105053e3b016174437ed83b25d2a79a88268f2")
	for {
		if len(hexBuf) == 0 {
			break
		}
		copy(tempBuf[:], hexBuf[0:32])
		testVals = append(testVals, tempBuf)
		hexBuf = hexBuf[32:]
	}
	result = TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "2d0ad2566627b50cd45125e89e963433b212b368cd2d91662c44813ba9ec90c2" {
		t.Fatal("Unable to hash a long object properly")
	}
}

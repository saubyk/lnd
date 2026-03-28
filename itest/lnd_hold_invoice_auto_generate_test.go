package itest

import (
	"context"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/lntest"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/stretchr/testify/require"
)

// testHoldInvoiceAutoGenerate tests that hold invoices can be created without
// providing a hash, with the server auto-generating a preimage and hash. It
// also tests the case where a preimage is provided and the server derives the
// hash, and verifies that both hash and preimage set returns an error.
func testHoldInvoiceAutoGenerate(ht *lntest.HarnessTest) {
	// Open a channel between Alice and Bob.
	_, nodes := ht.CreateSimpleNetwork(
		[][]string{nil, nil}, lntest.OpenChannelParams{
			Amt: btcutil.Amount(1_000_000),
		},
	)
	alice, bob := nodes[0], nodes[1]

	// Test 1: Auto-generate preimage and hash (neither provided).
	autoReq := &invoicesrpc.AddHoldInvoiceRequest{
		Memo:  "auto-generated",
		Value: 10_000,
	}
	autoResp := bob.RPC.AddHoldInvoice(autoReq)

	// The response must contain both a preimage and hash.
	require.Len(ht, autoResp.PaymentPreimage, 32,
		"expected 32-byte preimage")
	require.Len(ht, autoResp.PaymentHash, 32,
		"expected 32-byte payment hash")

	// Verify the hash matches SHA256 of the preimage.
	var autoPreimage lntypes.Preimage
	copy(autoPreimage[:], autoResp.PaymentPreimage)
	autoHash := autoPreimage.Hash()
	require.Equal(ht, autoHash[:], autoResp.PaymentHash,
		"hash should be SHA256 of preimage")

	// Subscribe, pay, and settle.
	autoStream := bob.RPC.SubscribeSingleInvoice(autoResp.PaymentHash)

	ht.SendPaymentAndAssertStatus(alice, &routerrpc.SendPaymentRequest{
		PaymentRequest: autoResp.PaymentRequest,
		FeeLimitSat:    1_000_000,
	}, lnrpc.Payment_IN_FLIGHT)

	ht.AssertInvoiceState(autoStream, lnrpc.Invoice_ACCEPTED)

	bob.RPC.SettleInvoice(autoResp.PaymentPreimage)
	ht.AssertInvoiceState(autoStream, lnrpc.Invoice_SETTLED)
	ht.AssertPaymentStatus(
		alice, autoHash, lnrpc.Payment_SUCCEEDED,
	)

	// Test 2: User-supplied preimage (server derives hash).
	var userPreimage lntypes.Preimage
	copy(userPreimage[:], ht.Random32Bytes())
	expectedHash := userPreimage.Hash()

	preimgResp := bob.RPC.AddHoldInvoice(
		&invoicesrpc.AddHoldInvoiceRequest{
			Memo:     "user-preimage",
			Value:    10_000,
			Preimage: userPreimage[:],
		},
	)

	// The response should echo back the preimage and correct hash.
	require.Equal(ht, userPreimage[:], preimgResp.PaymentPreimage,
		"preimage should be echoed back")
	require.Equal(ht, expectedHash[:], preimgResp.PaymentHash,
		"hash should match SHA256 of preimage")

	// Subscribe, pay, and settle.
	preimgStream := bob.RPC.SubscribeSingleInvoice(
		preimgResp.PaymentHash,
	)

	ht.SendPaymentAndAssertStatus(alice, &routerrpc.SendPaymentRequest{
		PaymentRequest: preimgResp.PaymentRequest,
		FeeLimitSat:    1_000_000,
	}, lnrpc.Payment_IN_FLIGHT)

	ht.AssertInvoiceState(preimgStream, lnrpc.Invoice_ACCEPTED)

	bob.RPC.SettleInvoice(userPreimage[:])
	ht.AssertInvoiceState(preimgStream, lnrpc.Invoice_SETTLED)
	ht.AssertPaymentStatus(
		alice, expectedHash, lnrpc.Payment_SUCCEEDED,
	)

	// Test 3: Traditional hash-only flow still works.
	var hashPreimage lntypes.Preimage
	copy(hashPreimage[:], ht.Random32Bytes())
	payHash := hashPreimage.Hash()

	hashResp := bob.RPC.AddHoldInvoice(
		&invoicesrpc.AddHoldInvoiceRequest{
			Memo:  "hash-only",
			Value: 10_000,
			Hash:  payHash[:],
		},
	)

	// Preimage should be empty since the server doesn't know it.
	require.Empty(ht, hashResp.PaymentPreimage,
		"preimage should be empty for hash-only")

	// Hash should match what we provided.
	require.Equal(ht, payHash[:], hashResp.PaymentHash,
		"returned hash should match provided hash")

	// Subscribe, pay, and settle.
	hashStream := bob.RPC.SubscribeSingleInvoice(payHash[:])

	ht.SendPaymentAndAssertStatus(alice, &routerrpc.SendPaymentRequest{
		PaymentRequest: hashResp.PaymentRequest,
		FeeLimitSat:    1_000_000,
	}, lnrpc.Payment_IN_FLIGHT)

	ht.AssertInvoiceState(hashStream, lnrpc.Invoice_ACCEPTED)

	bob.RPC.SettleInvoice(hashPreimage[:])
	ht.AssertInvoiceState(hashStream, lnrpc.Invoice_SETTLED)
	ht.AssertPaymentStatus(
		alice, payHash, lnrpc.Payment_SUCCEEDED,
	)

	// Test 4: Both hash and preimage should error.
	var bothPreimage lntypes.Preimage
	copy(bothPreimage[:], ht.Random32Bytes())
	bothHash := bothPreimage.Hash()

	_, err := bob.RPC.Invoice.AddHoldInvoice(
		context.Background(), &invoicesrpc.AddHoldInvoiceRequest{
			Memo:     "should-fail",
			Value:    10_000,
			Hash:     bothHash[:],
			Preimage: bothPreimage[:],
		},
	)
	require.Error(ht, err, "expected error when both hash "+
		"and preimage are set")
	require.Contains(ht, err.Error(),
		"cannot set both hash and preimage")
}

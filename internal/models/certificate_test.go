package models

import (
	"testing"
	"verisart-api/internal"
)

// Helpers and fixtures

func GetCertificateModel() *CertificateModel {
	userModel := NewUserModel()

	return NewCertificateModel(userModel)
}

var cert = &Certificate{
	Title:   "Blah",
	OwnerID: "ab12",
}

// Tests start here
func TestCreateCertificate(t *testing.T) {
	c := GetCertificateModel()

	err := c.CreateCertificate(cert)
	if err != nil {
		t.Errorf("An error occured trying to create certificate, %s", err)
	}

	expectedID := internal.GenerateIDFromString(cert.Title)
	if cert.ID != expectedID {
		t.Errorf("Expected certificate ID to be '%s', got '%s'", expectedID, cert.ID)
	}
}

func TestCreateCertificateAlreadyExists(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)
	err := c.CreateCertificate(cert)
	if err == nil {
		t.Errorf("Duplicate certificate was created")
	}
}

func TestUpdateCertificate(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)

	cert.Note = "Updated cert for dummy reasons"
	err := c.UpdateCertificate(cert)

	if err != nil {
		t.Errorf("Certificate did not update")
	}

	updatedCert, err := c.GetCertificateByID(cert.ID)
	if updatedCert.Note != cert.Note {
		t.Errorf("Certificate did not update accordingly, expected Note field to be %s, got %s", cert.Note, updatedCert.Note)
	}
}

func TestUpdateNonExistingCertificate(t *testing.T) {
	c := GetCertificateModel()

	err := c.UpdateCertificate(cert)
	if err == nil {
		t.Error("Excepted CertificateModel.UpdateCertificate to error out")
	}
}

func TestDeleteCertificate(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)
	c.DeleteCertificate(cert.ID)

	_, err := c.GetCertificateByID(cert.ID)
	if err == nil {
		t.Error("Certificate was not deleted")
	}
}

func TestGetCertificateByOwnerID(t *testing.T) {
	c := GetCertificateModel()

	certs := c.GetCertificatesByOwnerID("bob")
	if len(certs) != 0 {
		t.Errorf("Expected 0 certificate for owner ID 'bob', got %d", len(certs))
	}
	// Add two certificates for 'bob'
	cert.OwnerID = "bob"
	c.CreateCertificate(cert)

	otherCertForBob := &Certificate{
		Title:   "Blah2",
		OwnerID: "bob",
	}
	c.CreateCertificate(otherCertForBob)

	// Add one certificate for 'Notbob'
	otherCertNotForBob := &Certificate{
		Title:   "Blah3",
		OwnerID: "Notbob",
	}
	c.CreateCertificate(otherCertNotForBob)

	certs = c.GetCertificatesByOwnerID("bob")
	if len(certs) != 2 {
		t.Errorf("Expected 2 certificate for owner ID 'bob', got %d", len(certs))
	}

	notBobCerts := c.GetCertificatesByOwnerID("Notbob")
	if len(notBobCerts) != 1 {
		t.Errorf("Expected 1 certificate for owner ID 'Notbob', got %d", len(certs))
	}
}

func TestGetCertificateByID(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)

	_, err := c.GetCertificateByID(cert.ID)
	if err != nil {
		t.Error("Did not find certificate")
	}
}

func TestCheckCertificateOwnership(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)

	ok, _ := c.CheckCertificateOwnership(cert.ID, "abcddef")
	if ok {
		t.Error("Failed to assert invalid certificate ownership")
	}

	ok, _ = c.CheckCertificateOwnership(cert.ID, cert.OwnerID)
	if !ok {
		t.Error("Failed to assert valid certificate ownership")
	}

	_, err := c.CheckCertificateOwnership("nope", cert.OwnerID)
	if err == nil {
		t.Error("Expected an error when asserting ownership of a non-existant certificate")
	}
}

func TestCreateTransfer(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)

	tsx := Transfer{
		Email:  "bob@table.com",
		Status: Pending,
	}

	err := c.CreateTransfer(cert.ID, tsx)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateTransferNonExistingCert(t *testing.T) {
	c := GetCertificateModel()

	tsx := Transfer{}
	err := c.CreateTransfer("abasad", tsx)
	if err == nil {
		t.Error("Expected error when trying to create transfer on non existant certificate")
	}
}

func TestAcceptTransfer(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)

	tsx := Transfer{
		Email:  "bob@table.com",
		Status: Pending,
	}

	c.CreateTransfer(cert.ID, tsx)

	beforeID := cert.OwnerID
	c.AcceptTransfer(cert.ID)
	if beforeID == cert.OwnerID {
		t.Error("Transfer was not effective")
	}

}

func TestAcceptTranferNonExistingCert(t *testing.T) {
	c := GetCertificateModel()

	err := c.AcceptTransfer("sadasd")
	if err == nil {
		t.Error("Expected error got none")
	}
}

func TestAcceptTransferWhenNoneHasBeenCreated(t *testing.T) {
	c := GetCertificateModel()

	c.CreateCertificate(cert)

	err := c.AcceptTransfer(cert.ID)
	if err == nil {
		t.Error("Expected error got none")
	}

}

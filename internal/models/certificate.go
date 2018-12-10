package models

import (
	"fmt"
	"time"
	"verisart-api/internal"
)

// CertificateModel is used to interact with the certificate store
type CertificateModel struct {
	store     map[string]*Certificate
	userModel *UserModel
}

// Certificate represents a Verisart certificate
type Certificate struct {
	ID        string    `json:"id,omitempty"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	OwnerID   string    `json:"ownerId"`
	Year      int       `json:"year"`
	Note      string    `json:"note"`
	Transfer  Transfer
}

// A TransferStatus represents the status of a given certificate transaction
type TransferStatus string

const (
	// None indicates no transfer has been initiated
	None TransferStatus = "none"
	// Pending indicates the transfer has been initiated
	Pending TransferStatus = "pending"
)

// Transfer represents a transfer of certificate within Verisart
type Transfer struct {
	Email  string         `json:"email"`
	Status TransferStatus `json:"status,omitempty"`
}

// NewCertificateModel returns a certificate model
func NewCertificateModel(userModel *UserModel) *CertificateModel {
	model := &CertificateModel{
		store:     make(map[string]*Certificate),
		userModel: userModel,
	}
	return model
}

// CreateCertificate creates a new certificate
// Uniqueness is deemed to be on Certificate.Title, just because
func (c *CertificateModel) CreateCertificate(cert *Certificate) (err error) {
	id := internal.GenerateIDFromString(cert.Title)

	if _, ok := c.store[id]; ok {
		err = fmt.Errorf("Certificate already exists")
	} else {
		c.store[id] = cert
		cert.ID = id
	}

	return
}

// UpdateCertificate updates an existing certificate
// Returns an error if no matching certificate is found
func (c *CertificateModel) UpdateCertificate(cert *Certificate) (err error) {
	if _, ok := c.store[cert.ID]; !ok {
		err = fmt.Errorf("No certificate found with ID %s", cert.ID)
	} else {
		c.store[cert.ID] = cert
	}
	return
}

// DeleteCertificate deletes the given certificate ID, or does nothing
func (c *CertificateModel) DeleteCertificate(certID string) {
	delete(c.store, certID)
}

// GetCertificatesByOwnerID returns all certificates a user identified by ownerID owns
func (c *CertificateModel) GetCertificatesByOwnerID(ownerID string) (certificates []*Certificate) {
	// Due to our store being indexless this will be inefficient at any sort of scale
	// Effictively O(n) where n is the size of c.store
	for _, el := range c.store {
		if el.OwnerID == ownerID {
			certificates = append(certificates, el)
		}
	}

	if certificates == nil {
		certificates = make([]*Certificate, 0)
	}

	return
}

// GetCertificateByID returns the matching certificate if it exists
func (c *CertificateModel) GetCertificateByID(certID string) (cert *Certificate, err error) {
	var ok bool
	if cert, ok = c.store[certID]; !ok {
		err = fmt.Errorf("Certificacte not found")
	}
	return
}

// CheckCertificateOwnership asserts that a given certificateID is owned by the given ownerID
// error is only returned when the certificate is not found, for validation check owned only
func (c *CertificateModel) CheckCertificateOwnership(certID, ownerID string) (owned bool, err error) {
	owned = false
	if cert, ok := c.store[certID]; ok {
		if cert.OwnerID == ownerID {
			owned = true
		}
	} else {
		err = fmt.Errorf("Certificate does not exists")
	}
	return
}

// CreateTransfer marks a certificate as a pending transfer
func (c *CertificateModel) CreateTransfer(certID string, tsx Transfer) (err error) {
	if cert, ok := c.store[certID]; ok {
		cert.Transfer = tsx
	} else {
		err = fmt.Errorf("Certificate does not exists and thus cannot be transfered")
	}
	return
}

// AcceptTransfer transfers the certificate to the new owner
// if the email for the transfer is not attached to any existing account we create one
func (c *CertificateModel) AcceptTransfer(certID string) (err error) {
	cert, ok := c.store[certID]

	if !ok {
		err = fmt.Errorf("Certificate not found")
	} else {
		if cert.Transfer.Status != Pending {
			err = fmt.Errorf("No transfer has been initiated for this certificate")
		} else {
			owner, err := c.userModel.GetUserByEmail(cert.Transfer.Email)
			if err != nil {
				owner = &User{
					Name:  cert.Transfer.Email,
					Email: cert.Transfer.Email,
				}
				c.userModel.CreateUser(owner)
			}
			cert.OwnerID = owner.ID
			cert.Transfer = Transfer{
				Status: None,
			}
		}
	}

	return
}

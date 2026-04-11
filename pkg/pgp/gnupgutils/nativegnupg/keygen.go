package nativegnupg

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/pgp/gnupgutils/gnupgoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GenerateKeyPair(ctx context.Context, options *gnupgoptions.GenerateKeyPairOptions) (privateKey []byte, publicKey []byte, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	err = contextutils.CheckContextStillAlive(ctx)
	if err != nil {
		return nil, nil, err
	}

	logging.LogInfoByCtxf(ctx, "Generate GnuPG key pair started.")

	name, err := options.GetName()
	if err != nil {
		return nil, nil, err
	}

	comment, err := options.GetComment()
	if err != nil {
		return nil, nil, err
	}

	email, err := options.GetEmail()
	if err != nil {
		return nil, nil, err
	}

	rsaBits, err := options.GetRSABits()
	if err != nil {
		return nil, nil, err
	}

	logging.LogInfoByCtxf(ctx, "Going to generate new GnuPG key pair for '%s', comment='%s' email='%s' with %d RSABits", name, comment, email, rsaBits)

	entity, err := openpgp.NewEntity(name, comment, email, &packet.Config{
		RSABits: rsaBits,
	})
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to generate new GnuPG key: %w", err)
	}

	var pubKeyBuf bytes.Buffer
	pubKeyWriter, err := armor.Encode(&pubKeyBuf, "PGP PUBLIC KEY BLOCK", nil)
	if err != nil {
		tracederrors.TracedErrorf("Failed to create armor encoder: %w", err)
	}
	err = entity.Serialize(pubKeyWriter)
	if err != nil {
		tracederrors.TracedErrorf("Failed to serialize public key: %w", err)
	}
	pubKeyWriter.Close()

	publicKey = pubKeyBuf.Bytes()

	var privateKeyBuf bytes.Buffer
	privateKeyWriter, err := armor.Encode(&privateKeyBuf, "PGP PRIVATE KEY BLOCK", nil)
	if err != nil {
		tracederrors.TracedErrorf("Failed to create armor encoder: %w", err)
	}
	err = entity.SerializePrivate(privateKeyWriter, nil)
	if err != nil {
		tracederrors.TracedErrorf("Failed to serialize public key: %w", err)
	}
	privateKeyWriter.Close()

	privateKey = privateKeyBuf.Bytes()

	logging.LogInfoByCtxf(ctx, "Generate GnuPG key pair finished. Generated key has fingerprint '%s'.", getFingerprintFromEntity(entity))

	return privateKey, publicKey, nil
}

func getFingerprintFromEntity(entity *openpgp.Entity) string {
	return getFormatedFingerprint(entity.PrimaryKey.Fingerprint)
}

func getFormatedFingerprint(fingerprint []byte) string {
	fp := fmt.Sprintf("%X", fingerprint)
	formatted := fmt.Sprintf("%s %s %s %s %s  %s %s %s %s %s",
		fp[0:4], fp[4:8], fp[8:12], fp[12:16], fp[16:20],
		fp[20:24], fp[24:28], fp[28:32], fp[32:36], fp[36:40],
	)
	return formatted
}

func entityByPrivateKey(privateKey []byte) (entity *openpgp.Entity, err error) {
	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(privateKey))
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to read armored key ring: %w", err)
	}

	if len(entityList) == 0 {
		return nil, tracederrors.TracedErrorf("no keys found in armored key ring")
	}

	return entityList[0], nil
}

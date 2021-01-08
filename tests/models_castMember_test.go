package tests

import (
	"testing"
	"video-catalog/models"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfCastMemberIsEmpty(t *testing.T) {
	castMember := models.CastMember{}
	err := castMember.Validate()
	require.Error(t, err)
}

func TestValidateCastMemberFullFilled(t *testing.T) {
	castMember := models.CastMember{
		ID:   uuid.NewV4().String(),
		Name: "name",
		Type: models.TypeActor,
	}
	castMember.Prepare()
	err := castMember.Validate()
	require.Nil(t, err)
}

func TestValidateCastMemberIDisNotUUID(t *testing.T) {
	CastMember := models.CastMember{
		ID:   "id",
		Name: "name",
		Type: models.TypeActor,
	}
	CastMember.Prepare()
	err := CastMember.Validate()
	require.Error(t, err)
}

func TestCastMemberNameIsEmpty(t *testing.T) {
	castMember := models.CastMember{
		ID:   uuid.NewV4().String(),
		Type: models.TypeActor,
	}
	castMember.Prepare()
	err := castMember.Validate()
	require.Error(t, err)
}

func TestCastMemberNameIsLessThan3Characters(t *testing.T) {
	castMember := models.CastMember{
		ID:   uuid.NewV4().String(),
		Name: "ab",
		Type: models.TypeActor,
	}
	castMember.Prepare()
	err := castMember.Validate()
	require.Error(t, err)
}

func TestCastMemberNameIsMoreThan255Characters(t *testing.T) {
	castMember := models.CastMember{
		ID:   uuid.NewV4().String(),
		Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Type: models.TypeActor,
	}
	castMember.Prepare()
	err := castMember.Validate()
	require.Error(t, err)
}

func TestCastMemberTypeIsEmpty(t *testing.T) {
	castMember := models.CastMember{
		ID:   uuid.NewV4().String(),
		Name: "name",
	}
	castMember.Prepare()
	err := castMember.Validate()
	require.Error(t, err)
}

func TestCastMemberTypeWrongValue(t *testing.T) {
	castMember := models.CastMember{
		ID:   uuid.NewV4().String(),
		Name: "name",
		Type: 0,
	}
	castMember.Prepare()
	err := castMember.Validate()
	require.Error(t, err)
}

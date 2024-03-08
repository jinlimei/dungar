package learning

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadABook(t *testing.T) {
	assert.Equal(t, 1, 1)

	psg := `

If to be welcomed by the good,
O Book! thou make thy steady aim,
No empty chatterer will dare
To question or dispute thy claim.
But if perchance thou hast a mind
To win of idiots approbation,
Lost labour will be thy reward,
Though they’ll pretend appreciation.

They say a goodly shade he finds
Who shelters ’neath a goodly tree;
And such a one thy kindly star
In Bejar bath provided thee:
A royal tree whose spreading boughs
A show of princely fruit display;
A tree that bears a noble Duke,
The Alexander of his day.

Of a Manchegan gentleman
Thy purpose is to tell the story,
Relating how he lost his wits
O’er idle tales of love and glory,
Of “ladies, arms, and cavaliers:”
A new Orlando Furioso—
Innamorato, rather—who
Won Dulcinea del Toboso.

Put no vain emblems on thy shield;
All figures—that is bragging play.
A modest dedication make,
And give no scoffer room to say,
“What! Álvaro de Luna here?
Or is it Hannibal again?
Or does King Francis at Madrid
Once more of destiny complain?”

`
	//
	//	psg = `
	//
	//Of a Manchegan gentleman
	//Thy purpose is to tell the story,
	//Relating how he lost his wits
	//O’er idle tales of love and glory,
	//Of “ladies, arms, and cavaliers:”
	//A new Orlando Furioso—
	//Innamorato, rather—who
	//Won Dulcinea del Toboso.
	//
	//`

	val := ReadABook(psg)

	for pos, str := range val {
		fmt.Printf("%02d: %s\n\n", pos, str)
	}
}

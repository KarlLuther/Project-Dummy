package models

import (
	"errors"
)


/*the reason we have a custom error which we use instead of sql.ErrNoRows
The reason is to encapsulate data access layer logic and make it independent
from what kind of database we are using 
*/
var ErrNoRecord = errors.New("models: no matching record found")
package main

import (
	"fmt"
	"time"
)

func main() {
	// id, _ := uuid.NewV6()
	// fmt.Println(id)

	// var product1Id uuid.UUID = uuid.MustParse("018E8049-5A01-70F3-ACAE-4B6D5B5D91F5")
	// var product2Id uuid.UUID = uuid.MustParse("018E8049-5A25-7D3B-AD2B-677CC1AAF661")

	// ids := []uuid.UUID{product1Id, product2Id}
	// fmt.Println(ids)

	// sids := fmt.Sprintf("'%s'", strings.Join(uuid.UUIDs.Strings(ids), "','"))
	// fmt.Println(sids)

	fmt.Println(time.Now().UTC().Truncate(time.Second))
	fmt.Println(time.Now().UTC().Truncate(time.Millisecond))
	fmt.Println(time.Now().UTC().Truncate(time.Microsecond))
	fmt.Println(time.Now().UTC().Truncate(time.Nanosecond))

}

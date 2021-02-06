# conduit.go Output Validation

## Seemingly Incorrect Match 1)
```
"1": {
    "Name": "",
    "Path": "/Volumes/GoogleDrive/My Drive/Conduit_Tie_Ins/DYEA_LSA_8001729 (2).jpg",
    "ObjectId": 164361,
    "FQNID": "FIB:BUR::74001081",
    "WorkOrderName": "LSA_N_DATA_CORRECTION",
    "UDFs": [
        "INDUSTRY"
    ],
    "Point": {
        "lat": 33.99615,
        "lng": -117.96798611111112,
        "matched": false
    }
}
```
Okay, let's take a look at *DYEA_LSA_80001729 (2).jpg*
![DYEA_LSA_8001729 (2).jpg](images/ResultValidation/ImageLocation.png)

Now let's take a peak at where *2040 S. Hacienda Blvd* is in Google Maps
![Image Location by Address](images/ResultValidation/ImageLocByAddress.png)

Okay, now does this area match up in Live Maps?
![Permit Location Live Maps](images/ResultValidation/DyeaLocation.png)

That looks okay? What about the *FQNID*?
![FQNID Location Live Maps](images/ResultValidation/FqnidLocation.png)

Ahh, they're really close. I see why we got this FQNID.
The next question is, is this the right DYEA?

Let's check where the coordinates of the image put us -
![Image Location by Coords](images/ResultValidation/ImageLocByCoords.png)

Huh, I wonder what permit number comes up when I choose that corner -
![FQNID with Matching DYEA](images/ResultValidation/FQNIDwithDYEA.png)

**BINGO!**
So, the FQNID the program return *is* right after all.

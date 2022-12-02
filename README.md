# GOLANG-encryption-core
What ever im doing over in encryption-core I'm trying to better in GOLANG

WHAT IS THIS ?
* Your awnser is as good as mine. Eventually I want to make a gui file locker. In the meantime I assume the cli version needs to exist first.
  This is that. I origianally found a script that used encryption magic im learning about to encrypt files. I didn't know enough about bash
   scripting at the time so I made a seprate script to handle it and give me a more convient syntax for what I was working on. Now I'm just making it          stupidly complicated just messing around with new concepts I learn or want to try.

IS THIS COMPATIBLE WITH THE BASH VERSION ?
* No-ish* The concept is the same in go, python, and bash. But this version has a few diffrences.
  1. I finnaly realized how mush i was shooting myself in the foot with my variable names
  2. While learning about aes-cbc encryption the implementation I'm using currently only
    allows 32 byte keys. So while technically It would be possible to import the folders and 
    use diffrent versions to read and write the files. I havn't tryed it and it sounds like a bad idea
  3. On the note of variable names, The json names are diffrent. This will be eventually be remedyed in 
     the other versions but for now tomfoolery would definitly insue in you changed platforms somehow and 
     kept the same data.
  4. HMAC / crc the HMAC is 64bits appened to the end of of the file. when the file is read the HMAC is 
     removed and the HMAC is recalculated. I haven't tested how the original script works with HMAC vaalidation
     but it's probibly not my simple minded way.
  In conclusion. If your using this for anything important (why are you ?) this is not a good idea. if your just
     messing around like me go for it and try. The more we fuck around the more we find out !

The goal is to get a cycle going of bouncing between however many diffrent languages I port this too fixing adding features
     and keeping cross compatibility with them
  
 If you plan on doing something mission critical work I'd recommend using other software, the original script that inspired      this progect was really well written and well commented. I garunte nothing but an ever going quest to get better
 
HOW ARE YOU SUPPORTING THIS
      * Occasionally and naively. I'm new to github and proper progect management and not echo "here"ing my way through           debugging. If your invested in this progect hopefully you and pour a little knowlege into me and help make this a           cooler thing than it is.

package isc

import (
	"fmt"
	"io/ioutil"
)

func DefaultBanner1() {
	fmt.Printf("%s\n", defaultBanner)
}

func DefaultBanner2() {
	fmt.Printf("%s\n", BindDwenDwen)
}

func Banner(filePath string) {
	if b, err := ioutil.ReadFile(filePath); err == nil {
		str := string(b)
		fmt.Printf("%s\n", str)
	} else {
		fmt.Printf("%s\n", defaultBanner)
	}
}

var defaultBanner = `
               #########
              ############
              #############
             ##  ###########
            ###  ###### #####
            ### #######   ####
           ###  ########## ####
          ####  ########### ####
        #####   ###########  #####
       ######   ### ########   #####
       #####   ###   ########   ######
      ######   ###  ###########   ######
     ######   #### ##############  ######
    #######  ##################### #######
    #######  ##############################
   #######  ###### ################# #######
   #######  ###### ###### #########   ######
   #######    ##  ######   ######     ######
   #######        ######    #####     #####
    ######        #####     #####     ####
     #####        ####      #####     ###
      #####      ;###        ###      #
        ##       ####        ####

 :: iSysCore Service (GOLANG) ::

`

var BindDwenDwen = `
        .:          :j        
         i.,,......:Ei:       
        KWG        .,#        
        Wi .:t.  i:, ;,       
       .. ..jj;;;;fKL .       
         iD:,     .iLi    .   
        ,i.         .Df .LWLi 
      . .: ###    ###L:::KEW  
      ;,; W#f#W  E:###fiLWKW  
      .::##t#iL ; K###DttW# : 
     ..; W###W  # ,,##,if#W.  
     , :.#### K .j ###tj,#    
     , i. ##    ..  K#;:;.t   
    ...,:.      ..  .:f; ;    
   .#L..,f..   ....:,;,;      
   K#W..t:;.:...::,,:Li:      
   W#K:.,t;it;j;,iLjfj;,      
  EW#D;... :iii,..ji,f,       
  D##:......,:,;;,:,:,,       
  :EL  :...        .:,.    :: Bing Dwen Dwen ::   
   .   .:..    i;  .::     :: iSysCore Service (GOLANG) ::    
       .::....   ..:,         
        ::::..:;;;:,.         
         f::::::j,,;L         
        ,LLL:,,,;ifG.         
        fDKEL;.:GWWW.         
        tWWWK   #WW#i         
        ,j,j.   t#Kjt
`

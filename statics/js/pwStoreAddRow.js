function newPWStoreRow(elementId){
	var keyword = document.createElement("input");
    keyword.type = "text";
    keyword.name = "keyword";
    keyword.placeholder = "Keyword";
    keyword.className = "focus:outline-orange-500 focus:outline-2 text-base rounded-sm py-1 px-2 outline-gray-300 outline outline-1";
	var passphrase = document.createElement("input");
    passphrase.type = "text";
    passphrase.name = "value";
    passphrase.placeholder = "Passphrase";
    passphrase.className = "col-span-2 focus:outline-orange-500 focus:outline-2 text-base rounded-sm py-1 px-2 outline-gray-300 outline outline-1";

	document.getElementById(elementId).appendChild(keyword);
    document.getElementById(elementId).appendChild(passphrase);
    console.log("added")
   }

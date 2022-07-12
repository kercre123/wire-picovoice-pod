var eventListenerAdded = false

function updateIntentSelection() {
let xhr = new XMLHttpRequest();
xhr.open("GET", "/api/get_custom_intents_json");
xhr.setRequestHeader("Content-Type", "application/json");
xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
xhr.responseType = 'json';
xhr.send();
xhr.onload = function() {
  var listResponse = xhr.response
  if (listResponse != null) {
    
  var listNum = Object.keys(listResponse).length
  console.log(listNum)
  var select = document.createElement("select");
  document.getElementById("editSelect").innerHTML = ""
  select.name = "intents";
  select.id = "intents"
  for (const name in listResponse)
  {
    console.log(listResponse[name]["name"])
      var option = document.createElement("option");
      option.value = listResponse[name]["name"];
      option.text = listResponse[name]["name"]
      select.appendChild(option);
  }
  var label = document.createElement("label");
  label.innerHTML = "Choose the intent you would like to edit: "
  label.htmlFor = "intents";
  document.getElementById("editSelect").appendChild(label).appendChild(select);
} else {
    console.log("No intents found")
    var error = document.createElement("p");
    error.innerHTML = "No intents found, you must add one first"
    document.getElementById("editSelect").appendChild(error);
}
};
}

// get intent from editSelect element and create a form in div#editIntentForm to edit it
// options are: name, description, utterances, intent, paramname, paramvalue, exec
function editFormCreate() {
    var intentNumber = document.getElementById("intents").selectedIndex;
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/api/get_custom_intents_json");
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr.responseType = 'json';
    xhr.send();
    xhr.onload = function() {
        var intentResponse = xhr.response[intentNumber];
        if (intentResponse != null) {
            console.log(intentResponse)
            var form = document.createElement("form");
            form.id = "editIntentForm";
            form.name = "editIntentForm";
            var name = document.createElement("input");
            name.type = "text";
            name.name = "name";
            name.id = "name";
            // create label for name
            var nameLabel = document.createElement("label");
            nameLabel.innerHTML = "Name: "
            nameLabel.htmlFor = "name";
            name.value = intentResponse["name"];
            var description = document.createElement("input");
            description.type = "text";
            description.name = "description";
            description.id = "description";
            // create label for description
            var descriptionLabel = document.createElement("label");
            descriptionLabel.innerHTML = "Description: "
            descriptionLabel.htmlFor = "description";
            description.value = intentResponse["description"];
            var utterances = document.createElement("input");
            utterances.type = "text";
            utterances.name = "utterances";
            utterances.id = "utterances";
            // create label for utterances
            var utterancesLabel = document.createElement("label");
            utterancesLabel.innerHTML = "Utterances: "
            utterancesLabel.htmlFor = "utterances";
            utterances.value = intentResponse["utterances"];
            var intent = document.createElement("input");
            intent.type = "text";
            intent.name = "intent";
            intent.id = "intent";
            // create label for intent
            var intentLabel = document.createElement("label");
            intentLabel.innerHTML = "Intent: "
            intentLabel.htmlFor = "intent";
            intent.value = intentResponse["intent"];
            var paramname = document.createElement("input");
            paramname.type = "text";
            paramname.name = "paramname";
            paramname.id = "paramname";
            // create label for paramname
            var paramnameLabel = document.createElement("label");
            paramnameLabel.innerHTML = "Param Name: "
            paramnameLabel.htmlFor = "paramname";
            paramname.value = intentResponse["paramname"];
            var paramvalue = document.createElement("input");
            paramvalue.type = "text";
            paramvalue.name = "paramvalue";
            paramvalue.id = "paramvalue";
            // create label for paramvalue
            var paramvalueLabel = document.createElement("label");
            paramvalueLabel.innerHTML = "Param Value: "
            paramvalueLabel.htmlFor = "paramvalue";
            paramvalue.value = intentResponse["paramvalue"];
            var exec = document.createElement("input");
            exec.type = "text";
            exec.name = "exec";
            exec.id = "exec";
            // create label for exec
            var execLabel = document.createElement("label");
            execLabel.innerHTML = "Exec: "
            execLabel.htmlFor = "exec";
            exec.value = intentResponse["exec"];
            // create button that launches function
            var submit = document.createElement("button");
            submit.type = "button";
            submit.id = "submit";
            submit.innerHTML = "Submit";
            submit.onclick = function() {
                editIntent(intentNumber);
            }
            form.appendChild(nameLabel).appendChild(name);
            form.appendChild(document.createElement("br"));
            form.appendChild(descriptionLabel).appendChild(description);
            form.appendChild(document.createElement("br"));
            form.appendChild(utterancesLabel).appendChild(utterances);
            form.appendChild(document.createElement("br"));
            form.appendChild(intentLabel).appendChild(intent);
            form.appendChild(document.createElement("br"));
            form.appendChild(paramnameLabel).appendChild(paramname);
            form.appendChild(document.createElement("br"));
            form.appendChild(paramvalueLabel).appendChild(paramvalue);
            form.appendChild(document.createElement("br"));
            form.appendChild(execLabel).appendChild(exec);
            form.appendChild(document.createElement("br"));
            form.appendChild(submit);
            document.getElementById("editIntentForm").innerHTML = "";
            document.getElementById("editIntentForm").appendChild(form);
        } else {
            console.log("No intent found")
            var error = document.createElement("p");
            error.innerHTML = "No intent found, you must add one first"
            document.getElementById("editIntentForm").appendChild(error);
        }
    };
}

// create editIntent function that sends post to /api/edit_custom_intent, get index of intent to edit
// form data should include the intent number, name, description, utterances, intent, paramname, paramvalue, exec
function editIntent(intentNumber) {
    console.log(intentNumber)
        var formData = new FormData();
        formData.append("number", intentNumber+1);
        formData.append("name", document.getElementById("name").value);
        formData.append("description", document.getElementById("description").value);
        formData.append("utterances", document.getElementById("utterances").value);
        formData.append("intent", document.getElementById("intent").value);
        formData.append("paramname", document.getElementById("paramname").value);
        formData.append("paramvalue", document.getElementById("paramvalue").value);
        formData.append("exec", document.getElementById("exec").value);
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/edit_custom_intent");
        xhr.send(formData);
        xhr.onload = function() {
            var response = xhr.response;
            console.log(response);
                console.log("Intent edited")
                var success = document.createElement("p");
                success.innerHTML = "Intent edited"
                document.getElementById("editIntentStatus").appendChild(success);
        }
    
}

updateIntentSelection()

var HttpClient = function() {
    this.get = function(aUrl, aCallback) {
        var anHttpRequest = new XMLHttpRequest();
        anHttpRequest.onreadystatechange = function() { 
            if (anHttpRequest.readyState == 4 && anHttpRequest.status == 200)
                aCallback(anHttpRequest.responseText);
        }

        anHttpRequest.open( "GET", aUrl, true );            
        anHttpRequest.send( null );
    }
}

function sendIntentAdd() {
    const form = document.getElementById('intentAddForm');
    var data = "name=" + form.elements['name'].value + "&description=" + form.elements['description'].value + "&utterances=" + form.elements['utterances'].value + "&intent=" + form.elements['intent'].value + "&paramname=" + form.elements['paramname'].value + "&paramvalue=" + form.elements['paramvalue'].value + "&exec=" + form.elements['exec'].value;
    var client = new HttpClient();
    var result = document.getElementById('addIntentStatus');
    const resultP = document.createElement('p');
    resultP.textContent =  "Adding...";
    result.innerHTML = '';
    result.appendChild(resultP);
    fetch("/api/add_custom_intent?" + data)
    .then(response => response.text())
    .then((response) => {
        resultP.innerHTML = response
        result.innerHTML = '';
        result.appendChild(resultP);
        updateIntentSelection()
    })
}
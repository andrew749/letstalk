$(document).ready(function() {
    var createMessage = function(contents, messageClass) {
      return $("<div/>").addClass("alert").addClass(messageClass).text(contents);
    };

    var errorMessage = function() {
      return createMessage("Error, unable to create subscription", "alert-danger");
    };

    var successMessage = function() {
      return createMessage("Succesfully subscribed for updates!", "alert-success");
    };

    $("#signupBtn").click(function(e){
        $("#messageContainer").empty();
        if($("#signupForm")[0].checkValidity()) {
            var formData = $("#signupForm").serializeArray();
            var formObj = {};

            $.map(formData, function(n, i){
                if (n['name'] === 'classYear') {
                    formObj[n['name']] = parseInt(n['value']);
                } else {
                    formObj[n['name']] = n['value'];
                }
            });

            $.ajax({
                method: "POST",
                url: "https://api.hiveapp.org/v1/subscribe_email",
                headers: {
                    'Access-Control-Allow-Origin': '*'
                },
                data: JSON.stringify(formObj),
                success: function(response){
                  console.log(response);
                  if (response.status !== 200) {
                    var message = errorMessage();
                    $("#messageContainer").append(message);
                    return;
                  }
                  var message = successMessage();
                  $("#messageContainer").append(message);
                },
                error: function(obj){
                  console.log(obj);
                  var message = errorMessage();
                  $("#messageContainer").append(message);
                },
                dataType: "json",
                contentType : "application/json"
            });
            e.preventDefault();
        }
    });
});

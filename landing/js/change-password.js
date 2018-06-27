$(document).ready(function(){

  var baseUrl = "https://api.hiveapp.org/v1";

  function getPasswordChangeRequestId() {
    var searchParams = new URLSearchParams(window.location.search);
    return searchParams.get('requestId');
  }

  function submitPasswordChange(requestId, newPassword) {
    console.log("sending password change");
    $.post( baseUrl + "/change_password", JSON.stringify({
      "requestId": requestId,
      "password": newPassword,
    }), function( data ) {
      // on success change ui
      console.log("Successfully submitted change password request")
      passwordChangeHandler();
    });
  }

  var requestId = getPasswordChangeRequestId();

  $('form').submit(function(event) {
    event.preventDefault();
    console.log("got submit click");
    var newPassword = $('#newPassword').val();
    var newPasswordConfirm = $('#newPasswordConfirm').val();

    if (newPassword === "" || newPassword !== newPasswordConfirm) {
        badPasswordHandler();
        return;
    }

    submitPasswordChange(requestId, newPassword)
  });

  function badPasswordHandler() {
    console.log("Bad password");
    $("#message-container").html($('<div/>', {class: 'alert alert-danger'}).text("Failed to change password.").appendTo("#message-container"));
  }

  function passwordChangeHandler() {
    console.log("password changed");
    $("#message-container").html($('<div/>', {class: 'alert alert-success'}).text("Successfully changed password."));
  }
});

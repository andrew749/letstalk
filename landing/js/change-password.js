$(document).ready(function(){

  // let baseUrl = "https://api.hiveapp.org/";
  var baseUrl = "http://localhost:3000";

  function getPasswordChangeRequestId() {
    var searchParams = new URLSearchParams(window.location.search);
    return searchParams.get('requestId');
  }

  function submitPasswordChange(requestId, newPassword) {
    console.log("sending password change");
    $.post( baseUrl + "/change_password", {
      "requestId": requestId,
      "password": newPassword,
    }, function( data ) {
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
      // TODO: change styling of page to indicate bad password
  }

  function passwordChangeHandler() {
    console.log("password changed");
      // TODO: change styling of page to indicate a good password
  }
});

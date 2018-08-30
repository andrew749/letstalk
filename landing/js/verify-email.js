$(document).ready(function(){

  var baseUrl = "https://api.hiveapp.org/v1";

  function getVerifyEmailRequestId() {
    var searchParams = new URLSearchParams(window.location.search);
    return searchParams.get('requestId');
  }

  function setSuccessMessage(title, text) {
  }

  function setErrorMessage(title, text) {
  }

  function submitEmailVerifyRequest(requestId) {
    console.log("sending email verify request");
    $.post( baseUrl + "/verify_email", JSON.stringify({
      "requestId": requestId,
    })).done( function( data ) {
      // on success change ui
      console.log("Successfully submitted verify email request: ", data);
      setSuccessMessage("Email successfully verified")
    }).fail(function( data ) {
      console.log("Failed to verify email: ", data);
      setErrorMessage("Error verifying email")
    });
  }

  function handleRequest() {
    let requestId = getVerifyEmailRequestId();
    if (requestId === undefined || requestId === "") {
      setErrorMessage("Error: invalid email verification link");
      return;
    }
    submitEmailVerifyRequest();
  }

  handleRequest();

  function badPasswordHandler() {
    console.log("Bad password");
    $("#message-container").html($('<div/>', {class: 'alert alert-danger'}).text("Failed to change password.").appendTo("#message-container"));
  }

  function passwordChangeHandler() {
    console.log("password changed");
    $("#message-container").html($('<div/>', {class: 'alert alert-success'}).text("Successfully changed password."));
  }
});

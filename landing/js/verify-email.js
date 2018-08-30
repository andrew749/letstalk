$(document).ready(function(){

  var baseUrl = "http://localhost:80/v1";

  function getVerifyEmailRequestId() {
    var searchParams = new URLSearchParams(window.location.search);
    return searchParams.get('requestId');
  }

  function setSuccessMessage(msg) {
    $("#message-container").html($('<div/>', {class: 'alert alert-success'}).text(msg));
  }

  function setErrorMessage(msg) {
    $("#message-container").html($('<div/>', {class: 'alert alert-danger'}).text(msg).appendTo("#message-container"));
  }

  function submitEmailVerifyRequest(requestId) {
    console.log("sending email verify request for id: ", requestId);
    $.post( baseUrl + "/verify_email", JSON.stringify({
      "id": requestId,
    })).done( function( data ) {
      // on success change ui
      console.log("Successfully submitted verify email request: ", data);
      setSuccessMessage("Successfully verified email.")
    }).fail(function( data ) {
      console.log("Failed to verify email: ", data);
      setErrorMessage("Error verifying email. Please make sure to use the most recent email link you received.")
    });
  }

  function handleRequest() {
    var requestId = getVerifyEmailRequestId();
    if (!requestId) {
      console.error("missing request id in query string");
      setErrorMessage("Invalid email verification link.");
      return;
    }
    submitEmailVerifyRequest(requestId);
  }

  handleRequest();
});

$(document).ready(function(){

  var baseUrl = "https://api.hiveapp.org/v1";

  function getVerifyLinkRequestId() {
    var searchParams = new URLSearchParams(window.location.search);
    return searchParams.get('requestId');
  }

  function setSuccessMessage(msg) {
    $("#message-container").html($('<div/>', {class: 'alert alert-success'}).text(msg));
  }

  function setErrorMessage(msg) {
    $("#message-container").html(
      $('<div/>', {class: 'alert alert-danger'}).text(msg).appendTo("#message-container")
    );
  }

  function submitLinkVerifyRequest(requestId) {
    console.log("sending link verify request for id: ", requestId);
    $.post( baseUrl + "/verify_link", JSON.stringify({
      "verifyLinkId": requestId,
    })).done( function( data ) {
      // on success change ui
      console.log("Successfully submitted verify link request: ", data);
      setSuccessMessage("Successfully verified link.");
    }).fail(function( data ) {
      console.log("Failed to verify link: ", data);
      setErrorMessage(
        "Error verifying link. Please make sure to use the most recent link you received."
      );
    });
  }

  function handleRequest() {
    var requestId = getVerifyLinkRequestId();
    if (!requestId) {
      console.error("missing request id in query string");
      setErrorMessage("Invalid link verification link.");
      return;
    }
    submitLinkVerifyRequest(requestId);
  }

  handleRequest();
});

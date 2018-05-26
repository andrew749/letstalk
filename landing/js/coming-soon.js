$(function() {
  if (window.matchMedia('(min-width: 767px)').matches) {
    $(".img_highres").off().on("load", function() {
      var id = $(this).attr("id");
      var highres = $(this).attr("src").toString();
      var target = ".background";
      $(target).css("background-image", "url(" + highres + ")");
   });
  }
  
});
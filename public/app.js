'use strict';

var profile = "heap";
var tgdata = [];

$("#tgchart")
  .sparkline([1, 2, 4, 8, 2, 6, 1, 2, 3], {
    type: 'bar',
    height: '40px',
    barSpacing: 2,
    barColor: '#1ABC9C',
    barWidth: 5,
  });

refresh();

$(".cpu").on("click", function() {
  profile = "profile";
  refresh();
});
$(".heap").on("click", function() {
  profile = "heap";
  refresh();
});
$(".refresh").on("click", function() {
  refresh(true);
});
$(".filter").on("keyup", function() {
  refresh();
});

function refresh(force) {
  // TODO: cancel the existing request if it's not ciompleted.
  $('.results').html('Loading, be patient... CPU profile takes 30 seconds.');
  var f = $('.filter').val();
  $.get('/p', {
    profile: profile,
    filter: f,
    force: !!force
  }).done(function(data) {
    $('.results').html(data);
  }).fail(function(data) {
    $('.results').html(data);
  });
};
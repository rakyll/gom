'use strict';

(function(){
  var profile = "heap";
  var cumsort = true;
  var data = {
    max: 0,
    goroutine: [],
    thread: []
  };

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
  $('#cumsort').on('change', function() {
    cumsort = $(this).is(':checked');
    refresh();
  });

  var moreStats = function() {
    $.ajax('/stats').done(function(d) {
      appendChartData(data.goroutine, d.goroutine);
      appendChartData(data.thread, d.thread);
      drawCharts();
      setTimeout(function() {
        moreStats();
      }, 500);
    });
  };

  moreStats();

  function refresh(force) {
    // TODO: cancel the existing request if it's not ciompleted.
    $('.results').html('Loading, be patient... CPU profile takes 30 seconds.');
    var f = $('.filter').val();
    $.get('/p', { profile: profile, filter: f, cumsort: cumsort, force: !!force })
        .done(function(items) {
      var html = '';
      for (var i=0; i<items.length; i++) {
        var item = items[i];
        var row = '';
        row += '<td class="bar"><div style="width:' + item['score']*100 + 'px"></div></td>'
        row += '<td class="num">' + item['flat'] + '</td>';
        row += '<td class="num">' + item['flat_perc'] + '</td>';
        row += '<td class="num">' + item['flatsum_perc'] + '</td>';
        row += '<td class="num">' + item['cum'] + '</td>';
        row += '<td class="num">' + item['cum_perc'] + '</td>';
        row += '<td>' + item['name'] + '</td>';
        html += '<tr>' + row + '</tr>';
      }
      $(".results").html('<table>' + html + '</table>');
    }).fail(function(data) {
      $('.results').html(data);
    });
  };

  function appendChartData(target, val) {
    if (val > data.max) {
      data.max = val;
    }
    target.pop();
    if (target.length > 270) {
      target.shift();
    }
    target.push(val);
    // Add zeros, because sparkline draws the chart relatively,
    // depending on the min-max range of the dataset.
    target.push(0);
  }

  function drawCharts() {
    var opts = {
      type: 'line',
      height: '40px',
      lineColor: '#1ABC9C',
      lineWidth: 2,
      fillColor: '#e5e5e5',
      spotColor: '#1ABC9C',
      minSpotColor: '#1ABC9C',
      maxSpotColor: '#1ABC9C',
      chartRangeMax: data.max
    };
    $("#gorotinechart").sparkline(data.goroutine, opts);
    $("#threadchart").sparkline(data.thread, opts);
  }
})()

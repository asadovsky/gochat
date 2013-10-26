'use strict';

var msg = $('#msg');
var log = $('#log');

var appendToLog = function(msg) {
  var d = log[0];
  var doScroll = d && (d.scrollTop == d.scrollHeight - d.clientHeight);
  msg.appendTo(log);
  if (doScroll) {
    d.scrollTop = d.scrollHeight - d.clientHeight;
  }
};

var socket = new WebSocket($('body').attr('data-ws-url'));
socket.onclose = function(e) {
  appendToLog($('<div><strong>Socket closed.</strong></div>'));
};
socket.onmessage = function(e) {
  appendToLog($('<div></div>').text(e.data));
};

$('#form').submit(function() {
  if (!socket || !msg.val()) {
    return false;
  }
  socket.send(msg.val());
  msg.val('');
  return false;
});

$('#msg').focus();

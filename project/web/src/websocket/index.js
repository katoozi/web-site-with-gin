import ReconnectingWebSocket from "reconnecting-websocket";
import store, { addNotification } from "../redux/";
import MsgException from "./exception";

let hostname = window.location.hostname + ":" + window.location.port;
var socket = new ReconnectingWebSocket("ws://" + hostname + "/ws");

let connect = () => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    console.log("receive msg: ", msg);
    store.dispatch(addNotification(msg));
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

const sendMsg = async msg => {
  // const api = await fetch("https://httpbin.org/json");
  // try {
  //   const json = await api.json();
  //   console.log(json);
  // } catch (e) {
  //   console.log(e);
  //   throw new Error("we have error");
  // }
  try {
    socket.send(JSON.stringify(msg));
    return msg;
  } catch (e) {
    throw new MsgException("There is an error. Try again later!!!");
  }
};

const closeSocket = async (code, msg) => {
  socket.close(code, msg);
};

export { connect, sendMsg, closeSocket };
export { default as MsgException } from "./exception";

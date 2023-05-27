export default class GroundSegJS {

  constructor(url) {
    this.connected = false;
    this.url = url;
    this.structure = {};
    this.activity = {}
  }

  // Connect to websocket API
  connect() {
    console.log("attempting to connect..")
    this.ws = new WebSocket(this.url);
    return new Promise((resolve, reject) => {
      this.ws.onopen = () => {
        this.connected = true
        console.log("connected")
        resolve(this.connected)
      };

      this.ws.onmessage = (event) => {
        this.updateData(event.data)
      };

      this.ws.onerror = (error) => {
        console.log("connection failed", error);
      };

      this.ws.onclose = () => {
        this.connected = false
        console.log("closed")
      };
    });
  };

  // Login macro
  login(id,pwd="",token=null) {
    let data = {"id":id,"payload":{"category":"system","module":"login","action":{"password":pwd}}}
    console.log(id + " attempting to login.." )
    this.silentSend(data,token)
  }

  // Send token for verification
  verify(id,token=null) {
    let data = {"id":id,"payload":{"category":"token","module":null,"action":null}}
    console.log(id + " attempting to verify token.." )
    this.silentSend(data,token)
  }

  // Send raw action
  send(data,token=null) {
    if (token) {
      data['token'] = token
    }
    console.log(data.id + " attempting to send message.." )
    this.ws.send(JSON.stringify(data));
  }

  // Same as send but without logging
  silentSend(data,token=null) {
    if (token) {
      data['token'] = token
    }
    this.ws.send(JSON.stringify(data));
  }

  close() {
    this.ws.close();
  }

  deleteActivity(id) {
    if (id) {
      delete this.activity.activity[id]
      console.log(id + " activity acknowledged" )
    }
  }

  updateData(data) {
    console.log(data)
    /*
    data = JSON.parse(data)
    if (data.hasOwnProperty('activity')) {
      this.activity = this.deepMerge(this.activity, data)
    } else {
      this.structure = this.deepMerge(this.structure, data);
    }
    */
  };

  deepMerge(target, source) {
    for (const key in source) {
      if (typeof source[key] === 'object' && !Array.isArray(source[key]) && source[key] !== null) {
        if (!target.hasOwnProperty(key)) {
          target[key] = {};
        }
        this.deepMerge(target[key], source[key])
      } else {
        target[key] = source[key]
      }
    }
    return target
  }
}

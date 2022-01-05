'use strict';

class Options {
	onConnected = () => {};
	onDisconnected = () => {};
	onEvent = (method, params) => {};
}

class RPC {
	opts = new Options();

	ws = null;
	id = 0;
	pending = {};

	constructor(address, opts) {
		if (opts.onConnected) this.opts.onConnected = opts.onConnected;
		if (opts.onDisconnected) this.opts.onDisconnected = opts.onDisconnected;
		if (opts.onEvent) this.opts.onEvent = opts.onEvent;
		this.connect(address);
	}

	connect(address) {
		console.debug("Websocket: connecting to", address);
		let conn = new WebSocket(address);
		let _this = this;
		conn.onopen = () => {
			console.debug("Websocket: connected");
			_this.ws = conn;
			_this.opts.onConnected();
		};
		conn.onclose = () => {
			console.debug("Websocket: disconnected");
			_this.ws = null;
			_this.opts.onDisconnected();
			setTimeout(() => {
				_this.connect.call(_this, address);
			}, 3000);
		};
		conn.onmessage = (e) => {
			try {
				var d = JSON.parse(e.data);
				if (!d['id'] && d['method']) {
					console.debug("Websocket: event =>", e.data);
					_this.opts.onEvent(d.method, d.params);
					return
				}
				if (!!d['id'] && !d['method']) {
					if (!!_this.pending[d.id]) {
						console.debug("Websocket: response =>", e.data);
						if (d.hasOwnProperty('result')) _this.pending[d.id].resolve(d.result);
						else if (d.hasOwnProperty('error')) _this.pending[d.id].reject(d.error);
						delete _this.pending[d.id];
					}
					return
				}
			} catch (err) {
				console.error("Websocket: failed parsing message =>", err)
			}
		};
	}

	request(method, params) {
		let _this = this;
		return new Promise((resolve, reject) => {
			if (_this.ws == null) {
				reject("disconnected");
				return;
			}
			_this.id++;
			_this.pending[_this.id] = {
				resolve: resolve,
				reject: reject,
			};
			let req = JSON.stringify({ jsonrpc: "2.0", id: _this.id, method: method, params: params });
			console.trace("Websocket: request =>", req);
			_this.ws.send(req);
		});
	}
}

window.rpc = (address, opts) => {
	return new RPC(address, opts);
}
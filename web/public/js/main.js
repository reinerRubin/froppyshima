// burn it!


let gameState = {
  id: "",
  status: "init",
  requests: {},
  hits: {},

  canHit: function() {
    return this.status == "started"
  },
  setHit: function(x, y) {
    this.hits[`${x}${y}`] = 1;
  },
  getHit: function(x, y) {
    return this.hits[`${x}${y}`]
  },
  clear: function() {
    this.requests = {};
    this.hits = {};
  },
}

document.addEventListener("DOMContentLoaded", function() {
  // TODO mv to ENV
  start("ws://127.0.0.1:8080/ws")
});

function start(websocketServerLocation) {
  let webSocket = new WebSocket(websocketServerLocation);

  webSocket.onclose = function(){
    setTimeout(function(){
      start(websocketServerLocation)
    }, 5000);
  };

  webSocket.onmessage = function (event) {
    let data = JSON.parse(event.data);

    if (!data.ID && data.Body && data.Body.GameOver) {
      gameState.status = "finished"
      showMessage(data.Body.GameOver.Message)
    }

    if (!data.ID && data.Body && data.Body.ShipIsDestroyed) {
      showMessage(data.Body.ShipIsDestroyed.Message);
      showDestroyedShip(data.Body.ShipIsDestroyed.Ship);

    }

    if (!data.ID && data.Body && data.Body.ShipIsWounded) {
      showMessage(data.Body.ShipIsWounded.Message);
    }


    if (gameState.requests[data.ID] && gameState.requests[data.ID].method  == "hit") {
      let hitRequest = gameState.requests[data.ID],
	  x = hitRequest.x,
	  y = hitRequest.y,
	  cell = document.querySelectorAll(`[data-x="${x}"][data-y="${y}"]`)[0];

      gameState.setHit(x, y);

      if (data.Body.Result == 1) {
	cell.classList.add('killed');
      } else {
	cell.classList.add('checked');
      }
    }

    if (gameState.requests[data.ID] && gameState.requests[data.ID].method  == "new") {
      pureInit(10, 10);
      gameState.clear();

      gameState.id = data.Body.ID;
      showCode(data.Body.ID);

      bindCellClicks();
      gameState.status = "started";
    }

    if (gameState.requests[data.ID] && gameState.requests[data.ID].method  == "load") {
      if (!data.Error) {
	gameState.id = gameState.requests[data.ID].id;
	showCode(gameState.id);
	gameState.clear();
	load(data.Body);
	bindCellClicks();

	if (data.Body.AnyMoreShips) {
	  gameState.status = "started";
	} else {
	  gameState.status = "finished";
	}
      } else {
	showMessage(data.Error + ". A new game was started");
	newGame();
      }
    }
  }

  webSocket.onopen = function (event) {
    if (gameState.id) {
      let requestID = Math.random().toString();
      gameState.requests[requestID] = {
	method: "load",
	id: gameState.id,
      };

      gameState.status = "init";
      webSocket.send(JSON.stringify({ID: requestID, method: "load", body: {ID: gameState.id}}));

      return
    }

    // else new
    let requestID = Math.random().toString();
    gameState.requests[requestID] = {
      method: "new",
    }

    webSocket.send(JSON.stringify({ID: requestID, method: "new"}));
  };

  bindCellClicks = function() {
    let matrix = document.getElementById("matrix"),
	cells = matrix.getElementsByTagName("td");
    for (let i = 0; i < cells.length; i++) {
      cells[i].onclick =  function(event) {
	if (!gameState.canHit()) {
	  return
	}

	let x = +event.srcElement.dataset.x;
	let y = +event.srcElement.dataset.y;

	if (gameState.getHit(x, y)) return;

	let requestID = Math.random().toString();
	gameState.requests[requestID] = {
	  method: "hit",
	  x: x,
	  y: y,
	}

	webSocket.send(JSON.stringify({ID: requestID, method: "hit", body: {X: x, Y: y}}))
      }
    };
  }

  document.getElementById("load-by-code").onclick = function(event) {
    let id = document.getElementById('code-text').value;
    if (id == gameState.id) {
      return
    }

    let requestID = Math.random().toString();
    gameState.requests[requestID] = {
      method: "load",
      id: id,
    };

    gameState.status = "init";
    webSocket.send(JSON.stringify({ID: requestID, method: "load", body: {ID: id}}));
  }

  let newGame = function() {
    let requestID = Math.random().toString();
    gameState.requests[requestID] = {
      method: "new",
    };

    gameState.status = "init";
    webSocket.send(JSON.stringify({ID: requestID, method: "new", body: {}}));
  };

  document.getElementById("new").onclick = newGame;
}



function pureInit(maxx, maxy) {
  let table = document.getElementById("matrix");
  table.innerHTML = '';

  for (let y = 0; y < maxy; y++) {
    let row = table.insertRow(y);

    for (let x = 0; x < maxx; x++) {
      let cell = row.insertCell(x);
      cell.classList.add("cell");
      cell.dataset.x = x;
      cell.dataset.y = y;
    }
  }
}

function load(state) {
  pureInit(state.Maxx, state.Maxy);

  if (!state.AnyMoreShips) {
    showMessage("laded game was completed, please start a new one")
  }

  if (state.FailHits) {
    for (let i = 0; i < state.FailHits.length; i++) {
      let x = state.FailHits[i].X;
      let y = state.FailHits[i].Y;

      gameState.setHit(x, y);

      let cell = document.querySelectorAll(`[data-x="${x}"][data-y="${y}"]`)[0];
      cell.classList.add('checked');
    }
  }

  if (state.SuccessHits) {
    for (let i = 0; i < state.SuccessHits.length; i++) {
      let x = state.SuccessHits[i].X;
      let y = state.SuccessHits[i].Y;

      gameState.setHit(x, y);
      let cell = document.querySelectorAll(`[data-x="${x}"][data-y="${y}"]`)[0];
      cell.classList.add('killed');
    }
  }

  if (state.FatalHits) {
    for (let i = 0; i < state.FatalHits.length; i++) {
      let x0 = state.FatalHits[i].X;
      let y0 = state.FatalHits[i].Y;

      let cell = document.querySelectorAll(`[data-x="${x0}"][data-y="${y0}"]`)[0];
      cell.classList.add('killed');
    }

    let skirt = [-1, 0, 1]
    for (let i = 0; i < state.FatalHits.length; i++) {
      let x0 = state.FatalHits[i].X;
      let y0 = state.FatalHits[i].Y;

      for (let xsk = 0; xsk < skirt.length; xsk++) {
	for (let ysk = 0; ysk < skirt.length; ysk++) {
	  let x = x0 + skirt[xsk],
	      y = y0 + skirt[ysk];

	  gameState.setHit(x, y);

	  let cells = document.querySelectorAll(`[data-x="${x}"][data-y="${y}"]`);
	  if (cells[0]) {
	    if (!cells[0].classList.contains('killed')) {
	      cells[0].classList.add('checked');
	    }
	  }
	}
      }
    }
  }
}



function showDestroyedShip(ship) {
  let x0 = ship.Object.Coord.X,
      y0 = ship.Object.Coord.Y;

  for (let i = 0; i < ship.Segments.length; i++) {
    let xs = x0 + ship.Segments[i].Coord.X,
	ys = y0 + ship.Segments[i].Coord.Y;


    let skirt = [-1, 0, 1];

    for (let xsk = 0; xsk < skirt.length; xsk++) {
      for (let ysk = 0; ysk < skirt.length; ysk++) {
	let x = xs + skirt[xsk],
	    y = ys + skirt[ysk];

	let cells = document.querySelectorAll(`[data-x="${x}"][data-y="${y}"]`);
	if (cells[0]) {
	  if (!cells[0].classList.contains('killed')) {
	    gameState.setHit(x, y);
	    cells[0].classList.add('checked');
	  }
	}
      }
    }
  }
}

function showMessage(text) {
  let newMessage = document.createElement("div");
  newMessage.className = "message";
  newMessage.textContent = text;

  let messages = document.querySelector('.messages');
  messages.insertBefore(newMessage, messages.childNodes[0]);
}

function showCode(code) {
  let codeBox = document.getElementById('code-text');
  codeBox.value = code
}

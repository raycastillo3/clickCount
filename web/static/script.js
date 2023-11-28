window.addEventListener('load', function () {
    updateClickLabels()

})

//functionality
function itemClicked(item) {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function (){
        if (this.readyState == 4 && xhttp.status == 200){
            updateClickLabels();
            updateClickItem(),
            updateClickAddToCart(),
            updateClickBuy();
        }
    }
//PUT /api/clicks/{item} -> called when the {item}button was clicked once. no return data
    xhttp.open("PUT", `/api/clicks/${item}`, true)
    xhttp.send();
}

//asynchronous request
function updateClickLabels() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var response = JSON.parse(this.responseText);
            document.getElementById("item-label-item").innerHTML = "item: "  + response.itemClicks
            document.getElementById("item-label-addToCart").innerHTML = "addToCart: " + response.addToCartClicks
            document.getElementById("item-label-buy").innerHTML = "buy: " + response.buyClicks
            // console.log(response);
        }
    };
    xhttp.open("GET", "/api/clicks", true);
    xhttp.send();
}
//GET /api/clicks/item -> returns a single integer representing total item button clicks (`0`)
function updateClickItem() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function(){
        if (this.readyState == 4 && this.status == 200){
            var response = JSON.parse(this.responseText);
            document.querySelectorAll("item-button").innerHTML = "item: " + response.itemClicks
        }
    }
    xhttp.open("GET", "/api/clicks/item", true);
    xhttp.send();
}
//GET /api/clicks/addToCart -> returns a single integer representing total addToCart button clicks
function updateClickAddToCart() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function(){
        if (this.readyState == 4 && this.status == 200){
            var response = JSON.parse(this.responseText);
            document.querySelectorAll("item-button").innerHTML = "addToCart: " + response.addToCartClicks
        }
    }
    xhttp.open("GET", "/api/clicks/addToCart", true);
    xhttp.send();
}
//GET /api/clicks/buy -> returns a single integer representing total buy button clicks
function updateClickBuy() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function(){
        if (this.readyState == 4 && this.status == 200){
            var response = JSON.parse(this.responseText);
            document.querySelectorAll("item-button").innerHTML = "buy: " + response.buyClicks
        }
    }
    xhttp.open("GET", "/api/clicks/buy", true);
    xhttp.send();
}

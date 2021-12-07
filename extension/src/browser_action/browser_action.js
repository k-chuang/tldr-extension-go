var DEFAULT_TEXT = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.";
var GOOGLE_APP_ENGINE = "https://anchor-dot-i-amlg-dev.appspot.com/"

function processSelection(selection) {
  let e = chrome.runtime.lastError;
  if (e !== undefined) {
    // If there was an error getting the selection, we ignore it.
    // This will result in us sending a dummy DEFAULT_TEXT in the query.
    selection = null;
  }

  if (selection instanceof Array) {
    if (selection.length > 0) {
      // Valid Selection
      queryArticles(selection[0]);
    } else {
      // No Selection
      // TODO: Maybe use the whole page content instead?
      queryArticles(DEFAULT_TEXT);
    }
  } else {
    queryArticles(DEFAULT_TEXT);
  }
}

function queryArticles(query) {
  var request = { "query": query };

  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function () {
    if (xhr.readyState === XMLHttpRequest.DONE) {
      if (xhr.status === 200) {
        // We got a response from the server. The responseText is
        // a channel token so we can listen for a "verified" message.
        // token = xhr.responseText;
        // channel = new goog.appengine.Channel(token);
        // socket = channel.open(handler);
        var response = JSON.parse(xhr.response);
        renderResults(response["results"]);
      } else {
        //alert(`Bad Response from Server: ${xhr.status}`);
        // Fake response
        var response = JSON.parse(`{"results": {"Summary": "${DEFAULT_TEXT}", "Keywords": "None"}}`);
        // var response = JSON.parse(JSON.stringify(default_message));
        renderResults(response["results"]);
      }
    }
  }
  // xhr.open("POST", "/", true);
  xhr.open("POST", "http://localhost:8080/summarize", true);
  // xhr.open("POST", "https://anchor-dot-i-amlg-dev.appspot.com/", true);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.send(JSON.stringify(request));
}


function renderResults(results) {
  // var tableContent = "";
  var content = "";
  var summary = results["Summary"];
  var keywords = results["Keywords"]

  content += `
  <h2 class="mt-3">
    Summary
  </h2>
  <hr></hr>
  <p class="lead">
    ${summary}
  </p>
  <p class="h5"><strong>Keywords:</strong><em> ${keywords}</em></p>
  `;

  // `
  // <h3 class="mt-3">
  // Fancy display heading
  // <small class="text-muted">With faded secondary text</small>
  // </h3>
  // <hr></hr>
  // <p class="lead">
  //   Vivamus sagittis lacus vel augue laoreet rutrum faucibus dolor auctor. Duis mollis, est non commodo luctus.
  // </p>

  // <blockquote class="blockquote">
  // <p class="mb-0">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer posuere erat a ante.</p>
  // </blockquote>

  // `


  // tableContent += `
  // <table class="table table-striped">
  //   <colgroup>
  //     <col class="col-md-8">
  //   </colgroup>
  //   <thead>
  //     <tr>
  //       <th scope="col">Summary</th>
  //     </tr>
  //   </thead>
  //   <tbody id="results">
  //   <tr>
  //    <th scope = "col">${summary}%</th>
  //   </tr>
  // `;

  // for (var i=0; i < results.length; i++) {
  //   var entry_name = results[i]["name"];
  //   var entry_link = results[i]["link"];
  //   var entry_score = Math.round(parseFloat(results[i]["score"]) * 1000000) / 10000; // Round to 2 Decimal Places
  //   tableContent += `
  //   <tr>
  //     <th scope="col">
  //       <a href="${entry_link}" target="_blank">${entry_name}</a>
  //     </th>
  //     <td>${entry_score}%</td>
  //   </tr>
  //   `;
  // }
  // tableContent += `
  //   </tbody>
  // </table>
  // `;

  // document.getElementById("results").innerHTML = tableContent;
  document.getElementById("results").innerHTML = content;
}

// Try to get user selection for query
chrome.tabs.executeScript({
  code: "window.getSelection().toString();"
}, processSelection);

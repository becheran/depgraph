(function () {
  let cy = cytoscape({
    container: document.getElementById('cy'),

    layout: {
      // See: https://js.cytoscape.org/#layouts
      name: 'fcose',
      nodeDimensionsIncludeLabels: true,
      fit: true,
      idealEdgeLength: edge => 500,
    },

    style: [
      {
        "selector": "node",
        "style": {
          "width": 55,
          "height": 55,
          "text-valign": "center",
          'color': 'white',
          'text-outline-width': 1.5,
          'text-outline-color': '#888',
          'background-color': '#888'
        }
      },
      {
        "selector": "node[name]",
        "style": {
          "label": "data(name)"
        }
      },
      {
        "selector": "node[type = 'STD']",
        "style": {
          'background-color': 'red'
        }
      },
      {
        "selector": "node[type = 'PKG']",
        "style": {
          'background-color': 'blue'
        }
      },
      {
        "selector": "node[type = 'REQ']",
        "style": {
          'background-color': 'green'
        }
      },
      {
        "selector": "node[type = 'CGO']",
        "style": {
          'background-color': 'orange'
        }
      },
      {
        "selector": "edge",
        "style": {
          "width": 3,
          "target-arrow-shape": "triangle",
          "curve-style": "taxi",
          "control-point-step-size": 10
        }
      },
    ],

    elements: fetch('/api').then(response => response.json())
  });

  let detailsHeader = document.getElementById('details')

  cy.nodes().style()

  cy.on('mouseover', 'node', function (evt) {
    var node = evt.target;
    detailsHeader.innerHTML = node.id()
  });

  cy.on('mouseout', 'node', function (evt) {
    detailsHeader.textContent = ""
  });

  window.cy = cy
})();
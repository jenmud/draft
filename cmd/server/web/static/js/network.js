var network;

function convertJSON(data) {
    var options = {};
    var nodes = new vis.DataSet(options);
    var edges = new vis.DataSet(options);

    data.nodes.forEach(element => {
        var node = {
            id: element.uid,
            label: element.label,
            group: element.label,
            properties: {},
        }

        for (var key in element.properties) {
            node.properties[key] = atob(element.properties[key]);
        }

        console.debug(element);
        console.debug(node);
        nodes.add(node);
    });

    data.edges.forEach(element => {
        var edge = {
            id: element.uid,
            from: element.source_uid,
            label: element.label,
            to: element.target_uid,
            group: element.label,
            properties: {},
            arrows: "to",
        }

        for (var key in element.properties) {
            edge.properties[key] = atob(element.properties[key]);
        }

        console.debug(element);
        console.debug(edge);
        edges.add(edge);
    });

    return { "nodes": nodes, "edges": edges };
}

fetch("/assets/json")
    .then((resp) => { return resp.json(); })
    .then((dataJSON) => { return convertJSON(dataJSON); })
    .then((store) => {
        console.debug(store);

        var options = {
            height: '100%',
            width: '100%',
            nodes: {
                scaling: { min: 10, max: 20 },
                chosen: {
                    node: (values, id, selected, hovering) => {
                        values.color = "#ffe6e6";
                        values.shadow = true;
                    }
                },
            },
            edges: {
                chosen: {
                    edge: (values, id, selected, hovering) => {
                        values.color = "red";
                    },
                },
            },
            physics: {
                solver: "forceAtlas2Based",
            }
        };

        var container = document.getElementById('graph');
        var network = new vis.Network(container, store, options);

        network.on(
            "drag",
            () => {
                network.setOptions(
                    {
                        nodes: { physics: false },
                        edges: { physics: false },
                    }
                )
            }
        );
    })

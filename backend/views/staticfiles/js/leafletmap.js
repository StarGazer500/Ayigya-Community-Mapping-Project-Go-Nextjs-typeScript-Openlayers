// Global variables
var map, osm, layerControl, geoJsonLayerGroup;

// Function to initialize the map and add custom input control
function initializeMap() {
    // Create map
    map = L.map('map', {
        center: [6.2, -1.99],
        zoom: 10,
        zoomControl: false // Disable default zoom control (optional)
    });

    // Add OpenStreetMap tile layer
    osm = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</>'
    }).addTo(map);

    // Create layer control with base maps
    layerControl = L.control.layers({
        "OpenStreetMap": osm
    }, {}).addTo(map);

    // Add zoom control explicitly (if it's disabled or removed)
    L.control.zoom({
        position: 'topleft' // Positioning it in the top-left
    }).addTo(map);
    


}

// Function to create the input element



function addGeoJsonLayer1(data) {
    // Remove existing GeoJSON layers from the map and layer control
    if (geoJsonLayerGroup) {
        map.removeLayer(geoJsonLayerGroup);
        layerControl.removeLayer(geoJsonLayerGroup);
    }
    
    // Recreate the feature group
    geoJsonLayerGroup = L.featureGroup();
    
    // Function to convert snake_case to Title Case
    function formatKey(key) {
        return key
            .split('_')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ');
    }
    
    // Function to format value based on key
    function formatValue(key, value) {
        // Handle specific formatting for known keys
        const formatMap = {
            'shape__len': `${value} meters`,
            'shape__are': `${value} sq meters`,
            'creationda': value ? new Date(value).toLocaleDateString() : value
        };
        
        return formatMap[key] || value;
    }
    
    // Add new GeoJSON layers
    for(var i = 0; i < data.length; i++){
        var geoJsonLayer = L.geoJSON(data[i].geom, {
            style: function(feature) {
                return {
                    color: 'blue',
                    fillColor: 'blue',
                    fillOpacity: 0.3
                };
            },
            onEachFeature: function(feature, layer) {
                // Dynamically generate popup content
                var popupContent = '<div>';
                
                // Get all keys from the current data object
                Object.keys(data[i]).forEach(key => {
                    // Skip 'geom' key and keys with null/undefined/empty values
                    if (key !== 'geom' && data[i][key] != null && data[i][key] !== '') {
                        popupContent += `
                            <strong>${formatKey(key)}:</strong> ${formatValue(key, data[i][key])}<br>
                        `;
                    }
                });
                
                popupContent += '</div>';
                
                // Bind popup
                layer.bindPopup(popupContent);
                
                // Optional: Add click event
                layer.on('click', function(e) {
                    console.log('Feature clicked:', data[i]);
                });
            }
        });
        
        // Add layer to the feature group
        geoJsonLayerGroup.addLayer(geoJsonLayer);
    }
    
    // Add feature group to the map
    geoJsonLayerGroup.addTo(map);
    
    // Add to layer control
    layerControl.addOverlay(geoJsonLayerGroup, "GeoJSON Layers");
    
    // Zoom to layers if any exist
    if (geoJsonLayerGroup.getLayers().length > 0) {
        map.fitBounds(geoJsonLayerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
        });
    }
}

function addGeoJsonLayer(data) {
    // Check if DataTable exists and destroy it if it does
    if ($.fn.DataTable.isDataTable('#dataTable')) {
        $('#dataTable').DataTable().destroy();
    }
    
    // Diagnostic logging
    console.log('addGeoJsonLayer called with data:', data);
    
    // Check for critical DOM elements
    const tableHeader = document.getElementById('tableHeader');
    const tableBody = document.getElementById('tableBody');
    
    if (!tableHeader || !tableBody) {
        console.error('Table header or body not found. Check your HTML structure.');
        return;
    }
    
    // Check for map, layer control, and layer group
    if (!map || !layerControl) {
        console.error('Map or LayerControl not initialized');
        return;
    }
    
    // Check if data is empty
    if (!data || data.length === 0) {
        console.warn('No data provided to addGeoJsonLayer');
        tableHeader.innerHTML = '';
        tableBody.innerHTML = '';
        return;
    }
    
    // Remove existing GeoJSON layers from the map and layer control
    if (geoJsonLayerGroup) {
        map.removeLayer(geoJsonLayerGroup);
        layerControl.removeLayer(geoJsonLayerGroup);
    }
    
    // Recreate the feature group
    geoJsonLayerGroup = L.featureGroup();
    
    // Clear existing table
    tableHeader.innerHTML = '';
    tableBody.innerHTML = '';
    
    // Function to convert snake_case to Title Case
    function formatKey(key) {
        return key
            .split('_')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1))
            .join(' ');
    }
    
    // Function to format value based on key
    function formatValue(key, value) {
        // Handle specific formatting for known keys
        const formatMap = {
            'shape__len': `${value} meters`,
            'shape__are': `${value} sq meters`,
            'creationda': value ? new Date(value).toLocaleDateString() : value
        };
        return formatMap[key] || value;
    }
    
    // Prepare table headers dynamically
    const headers = Object.keys(data[0]).filter(key => key !== 'geom');
    console.log('Table headers:', headers);
    
    headers.forEach(header => {
        const th = document.createElement('th');
        th.textContent = formatKey(header);
        tableHeader.appendChild(th);
    });
    
    // Prepare data for DataTables
    const tableData = data.map((item, index) => {
        const rowData = {};
        headers.forEach(key => {
            rowData[formatKey(key)] = formatValue(key, item[key]);
        });
        rowData.index = index; // Store index for map interaction
        return rowData;
    });

    $('#map').css('height', '50vh');
    $('.leaflet-control-results-control-input').css('bottom', '20px');
     

    
    // Add new GeoJSON layers
    for(var i = 0; i < data.length; i++){
        // Create GeoJSON layer
        var geoJsonLayer = L.geoJSON(data[i].geom, {
            style: function(feature) {
                return {
                    color: 'blue',
                    fillColor: 'blue',
                    fillOpacity: 0.3
                };
            },
            onEachFeature: function(feature, layer) {
                // Dynamically generate popup content
                var popupContent = '<div>';
                Object.keys(data[i]).forEach(key => {
                    if (key !== 'geom' && data[i][key] != null && data[i][key] !== '') {
                        popupContent += `
                            <strong>${formatKey(key)}:</strong> ${formatValue(key, data[i][key])}<br>
                        `;
                    }
                });
                popupContent += '</div>';
                
                // Bind popup
                layer.bindPopup(popupContent);
                
                // Modify layer click to highlight corresponding row
                layer.on('click', function(e) {
                    // Trigger row selection in DataTable
                    const dataTable = $('#dataTable').DataTable();
                    const rowIndex = dataTable.rows().indexes().filter(function(idx) {
                        return dataTable.row(idx).data().index === i;
                    });
                    
                    if (rowIndex.length) {
                        // Select the row
                        dataTable.row(rowIndex[0]).select();
                        
                        // Zoom to layer
                        map.fitBounds(layer.getBounds(), {
                            padding: [50, 50],
                            maxZoom: 15
                        });
                        
                        // Open popup
                        layer.openPopup();
                        
                        // Scroll to row
                        const rowNode = dataTable.row(rowIndex[0]).node();
                        rowNode.scrollIntoView({
                            behavior: 'smooth',
                            block: 'center'
                        });
                    }
                });
            }
        });
        
        // Add layer to the feature group
        geoJsonLayerGroup.addLayer(geoJsonLayer);
    }
    
    // Add feature group to the map
    geoJsonLayerGroup.addTo(map);
    
    // Add to layer control
    layerControl.addOverlay(geoJsonLayerGroup, "GeoJSON Layers");
    
    // Zoom to layers if any exist
    if (geoJsonLayerGroup.getLayers().length > 0) {
        map.fitBounds(geoJsonLayerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
        });
    }
    
    // Initialize DataTable with additional features
    const dataTable = $('#dataTable').DataTable({
        data: tableData,
        columns: [
            ...headers.map(header => ({
                data: formatKey(header),
                title: formatKey(header)
            }))
        ],
        select: {
            style: 'single'
        },
        // Optional: Add more DataTables configurations as needed
        responsive: true,
        pageLength: 10,
        lengthMenu: [[10, 25, 50, -1], [10, 25, 50, "All"]]
    });
    
    // Add click event to table rows to interact with map
    $('#dataTable tbody').on('click', 'tr', function() {
        const data = dataTable.row(this).data();
        const index = data.index;
        
        // Get corresponding layer
        const layerGroup = geoJsonLayerGroup.getLayers()[index];
        
        if (layerGroup) {
            // Zoom to layer
            map.fitBounds(layerGroup.getBounds(), {
                padding: [50, 50],
                maxZoom: 15
            });
            
            // Open popup
            layerGroup.openPopup();
        }
    });
    
    console.log('DataTable initialized with rows:', tableData.length);
}


initializeMap();

L.Control.SearchInput = L.Control.extend({
    options: {
        position: 'topleft' // Position of the control
    },

    onAdd: function (map) {
        const container = L.DomUtil.create('div', 'leaflet-control leaflet-control-search-input');

        const input = L.DomUtil.create('input', 'search-input', container);
        input.placeholder = 'Search here...';

        // Add event listener for keyup event to trigger on Enter key press
        input.addEventListener('keyup', function(event) {
            const searchTerm = event.target.value;

            // Check if the pressed key is Enter (key code 13)
            if (event.key === 'Enter') {
                console.log("Search term:", searchTerm);
                
                // Only send the request if there's a search term
                if (searchTerm.trim() !== '') {
                    // Example request
                    $.ajax({
                        url: '/map/simplesearch', // Your server-side URL
                        type: 'POST',
                        dataType: 'json', // Expecting JSON response
                        contentType: 'application/json',
                        data: JSON.stringify({ searchValue: searchTerm }),
                        success: function (response) {
                            if (response.success) {
                                console.log("Successfully received data:", response.data);
                                // Process response here (e.g., adding GeoJSON layers)
                                var data =  response.data
                                $(".results-returned").val(data.length)
                                addGeoJsonLayer1(response.data);
                            } else {
                                alert("No results found.");
                            }
                        },
                        error: function (xhr, status, error) {
                            var errorMessage = xhr.responseJSON ? xhr.responseJSON.error : xhr.responseText || 'An unknown error occurred.';
                            alert("Error: " + errorMessage);
                        }
                    });
                } else {
                    alert("Please enter a valid search term.");
                }
            }
        });

        return container;
    }
});


// Add the control to the map
map.addControl(new L.Control.SearchInput());


L.Control.ResultsReturned = L.Control.extend({
    options: {
        position: 'bottomleft' // Position of the control
    },

    onAdd: function (map) {
        const container = L.DomUtil.create('div', 'leaflet-control leaflet-control-results-control-input');

        const input = L.DomUtil.create('input', 'results-returned', container);
        input.placeholder = 'no. of items appears here...';

        // Add event listener for search functionality
        input.addEventListener('keyup', function(event) {
            const searchTerm = event.target.value;

            
          
        });

        return container;
    }
});

// Add the control to the map
map.addControl(new L.Control.ResultsReturned());
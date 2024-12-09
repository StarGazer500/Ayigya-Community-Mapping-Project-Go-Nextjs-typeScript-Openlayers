import React, { useState,useEffect } from 'react';
import L from 'leaflet';
import { Map as LeafletMap } from "leaflet"; 


function SimpleQuery({addLayerToMap, setSimpleQueryResults,setIsTableViewOpen,map,setTableColumnNames,seTableData }) {
    const [value, setValue] = useState('');
    const [layerColor, setLayerColor] = useState('rgba(255, 0, 0, 1)'); 
    const [geojsonLayer, setGeojsonLayer] = useState<any>(null);
    const [bounds,setBounds] = useState()
   
    const serverbaseurl = "http://localhost:8080"

    const handleValueChange = (event) => {
        setValue(event.target.value);
    };

    useEffect(() => {
        if (bounds) {
            // Make sure bounds are valid before calling fitBounds
            map.fitBounds(bounds);
        }
      }, [bounds])

    const baseUrl = 'http://localhost:8080/geoserver/nurc/ows';

    const constructSearchQuery = async (searchValue:any):Promise<any> => {

        // async function getDataByColumn(path: string,table:string,column:string, map: LeafletMap): Promise<any> {
            try {
              // Fetch GeoJSON data from the server
              const response = await fetch(`${serverbaseurl}/map/simplesearch`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify({searchValue:searchValue}) // Send selectedLayer in the body
              });
          
              if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
              }
          
              let data = await response.json();
              data= data.data
              console.log("data",data)
          
              // Ensure data is an array
              if (!Array.isArray(data) || data.length === 0) {
                console.log("No data found or data is not an array.");
                return []; // Return an empty array if no valid data
              }

              setSimpleQueryResults(data.length)
          
              if (geojsonLayer) {
                map.removeLayer(geojsonLayer);
                // layerControl.removeLayer(geoJsonLayerGroup);
            }
          
              // Initialize a Leaflet feature group to hold the GeoJSON layers
              const geoJsonLayerGroup = L.featureGroup();
          
              // Function to format key (makes snake_case keys more human-readable)
              function formatKey(key: string): string {
                return key
                  .split('_')
                  .map(word => word.charAt(0).toUpperCase() + word.slice(1))
                  .join(' ');
              }
          
              // Function to format value based on key
              function formatValue(key: string, value: any): any {
                const formatMap = {
                  'shape__len': `${value} meters`,
                  'shape__are': `${value} sq meters`,
                  'creationda': value ? new Date(value).toLocaleDateString() : value
                };
                return formatMap[key] || value;
              }
          
              // Iterate over the data to create GeoJSON layers
              data.forEach((item: any) => {
                // Ensure item.geom exists and is valid
                if (item.geom) {
                  const geoJsonLayer = L.geoJSON(item.geom, {
                    style: function (feature) {
                      return { color: 'red' }; // You can customize the style here
                    },
                    onEachFeature: function (feature, layer) {
                      // Dynamically generate popup content
                      let popupContent = '<div>';
          
                      // Get all keys from the current data object
                      Object.keys(item).forEach(key => {
                        // Skip 'geom' key and keys with null/undefined/empty values
                        if (key !== 'geom' && item[key] != null && item[key] !== '') {
                          popupContent += `
                            <strong>${formatKey(key)}:</strong> ${formatValue(key, item[key])}<br>
                          `;
                        }
                      });
          
                      popupContent += '</div>';
          
                      // Bind popup
                      layer.bindPopup(popupContent);
          
                      // Optional: Add click event
                      layer.on('click', function (e) {
                        console.log('Feature clicked:', item);
                      });
                    }
                  });
          
                  // Add the created GeoJSON layer to the feature group
                  geoJsonLayerGroup.addLayer(geoJsonLayer);
                  
                } else {
                  console.warn("Missing 'geom' data for item:", item);
                }
              });
          
              // Add the feature group to the map
              setGeojsonLayer(geoJsonLayerGroup)
              geoJsonLayerGroup.addTo(map);
          
              // Optionally, zoom to the bounds of the GeoJSON layer
              map.fitBounds(geoJsonLayerGroup.getBounds(), { padding: [50, 50], maxZoom: 15 });
          
              console.log("GeoJSON data retrieved and layers added:", data.length);
          
              // Return the GeoJSON layer group, or the data if needed
              return geoJsonLayerGroup;
          
            } catch (error) {
              console.log("Error fetching GeoJSON data:", error);
              return []; // Return an empty array in case of error
            } finally {
              console.log("Map projection:", 'EPSG:3857'); // Leaflet uses Web Mercator by default
            }
        //   }

     
    };

    const handleKeyDown = async(event) => {
        if (event.key === 'Enter') {
            await constructSearchQuery(value)
        }
    };

    return (
        <input 
            id="results" 
            value={value}
            onChange={handleValueChange}
            onKeyDown={handleKeyDown}
            style={{ background: "brown" }} 
            placeholder='Search Anything....'
            className="absolute top-2 z-[5000] left-80 p-1 placeholder:text-white placeholder:opacity-100 rounded-md text-md text-white font-bold"
        />
    );
}

export default SimpleQuery;
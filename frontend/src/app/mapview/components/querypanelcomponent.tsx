import React, { useState,useCallback, useEffect} from 'react';
import L from 'leaflet';
import { Map as LeafletMap } from "leaflet"; 
// import WMSCapabilities from 'wms-capabilities'


function QueryPanel({addLayerToMap ,setIsTableViewOpen,results,setResults,tableData,setTableColumnNames,tableColumnNames,seTableData,isOpen,map}) {
  const [selectedAttribute, setSelectedAttribute] = useState('None');
  const [selectedFeatureLayer, setSelectedFeatureLayer] = useState('All Feature Layers');
  const [selectedQueryType, setSelectedQueryType] = useState('None');
  const [value, setValue] = useState('');
  // const [results, setResults] = useState('');
  const [layers, setLayers] = useState<unknown[]>([]);
  const [attributes,setAttributes] = useState<unknown[]>([]);
  const [operators,setOperations]= useState<unknown[]>([]);
  const [error, setError] = useState(null);
  const [geojsonLayer, setGeojsonLayer] = useState<any>(null);
  const [bounds,setBounds] = useState()
  
  const serverbaseurl = "http://localhost:8080"
 

  const handleSelectedAttributeOnChange = (event) => {
    setSelectedAttribute(event.target.value);
  };

  const handleSelectedFeatureLayerOnChange = (event) => {
    setSelectedFeatureLayer(event.target.value);
    
  };

  const handleSelectedQueryTypeOnChange = (event) => {
    setSelectedQueryType(event.target.value);
  };

  const handleValueChange = (event) => {
    setValue(event.target.value);
  };

  const handleResultsChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setResults(event.target.value);
  };


 

  async function getAttributeNames(layer:any) {
   
     
  
      fetch(`${serverbaseurl}/map/featureattributes`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ selectedLayer:layer}) // Send searchLayer in the body
      })
        .then(response_data => {
          if (!response_data.ok) {
            console.log(response_data.status);
            throw new Error('Network response was not ok');
          }
          return response_data.json();
        })
        .then(response_data => {
          // console.log("returned data", response_data.data);
          const newLayers = Object.values(response_data.data);
          setAttributes(newLayers);
          console.log("attributes",response_data.data );
        })
        .catch(error => {
          console.log('Unknown error occurred:', error.message);
        });
   
   
   }



// async function getData(layerName:any, map:LeafletMap):Promise<any> {


//   try {
//     const response = await fetch(`${serverbaseurl}/map/featureoperatures`, {
//             method: 'POST',
//             credentials: 'include',
//             headers: {
//               "Content-Type": "application/json",
//             },
//             body: JSON.stringify({ sselectedLayer:layerName}) // Send searchLayer in the body
//           })
//     if (!response.ok) {
//       throw new Error(`HTTP error! status: ${response.status}`);
//     }
//     const data = await response.json();
//     console.log("GeoJSON data received:", data);

//     // Create a Leaflet GeoJSON layer
//     for  (var i=0;i<data.length;i++){
//     const geoJsonLayer = L.geoJSON(data[i].geom, {
//       style: function (feature) {
//         return { color: 'blue' }; // You can customize the style here
//       },
//       // onEachFeature: function (feature, layer) {
//       //   // You can add popups or tooltips here
//       //   if (feature.properties) {
//       //     layer.bindPopup(Object.keys(feature.properties).map(key => 
//       //       feature.properties[key]!==null?`${key}: ${feature.properties[key]}`:null
//       //     ).join('<br>'));
//       //   }
//       // }
//     });
//   }

   

//   //   // Fit the map to the GeoJSON layer's bounds
//   //   const bounds = geoJsonLayer.getBounds()
//   //   const latLngBounds = [
//   //     [bounds._southWest.lat, bounds._southWest.lng], // Southwest corner
//   //     [bounds._northEast.lat, bounds._northEast.lng]  // Northeast corner
//   // ]

//   // setBounds(latLngBounds)

//   // console.log("bounds",latLngBounds)

//   //   console.log("Features loaded:", data.features.length);
//   //   setResults(data.features.length);

//     if (data.features.length > 0) {
//       // Update tableData with the new features
//       seTableData(data.features);
//       console.log("Data Retrieved Successfully",);
//       console.log("First feature geometry:", data.features[0].geometry.coordinates);
//       return data.features; // Return the features if any are found
//     } else {
//       console.log("No features found for the query.");
//       return []; // Return an empty array if no features are found
//     }
//   } catch (error) {
//     console.error("Error fetching GeoJSON data:", error);
//     return null; // Return null or handle the error as needed
//   } finally {
//     console.log("Map projection:", 'EPSG:3857'); // Leaflet uses Web Mercator by default
//   }
// }

async function getAllData(path: string, map: LeafletMap): Promise<any> {
  try {
    // Fetch GeoJSON data from the server
    const response = await fetch(`${serverbaseurl}/map/${path}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({}) // Send selectedLayer in the body
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
    setResults(data.length)

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
          title: item.TableName ,
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
    console.error("Error fetching GeoJSON data:", error);
    return []; // Return an empty array in case of error
  } finally {
    console.log("Map projection:", 'EPSG:3857'); // Leaflet uses Web Mercator by default
  }
}


async function getDataByColumn(path: string,table:string,column:string, map: LeafletMap): Promise<any> {
  try {
    // Fetch GeoJSON data from the server
    const response = await fetch(`${serverbaseurl}/map/${path}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({selectedLayer:table,selectedAttribute:column}) // Send selectedLayer in the body
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

    setResults(data.length)
    
    seTableData(data);
    setIsTableViewOpen(true)

    
    

    // Set table column names from the first record (you can adjust this logic)
    const columns = Object.keys(data[0]).map((key) => ({
      name: formatKey(key)
    }));

    setTableColumnNames(columns);

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
    // console.log("selectedlayer",selectedFeatureLayer)

    // Iterate over the data to create GeoJSON layers
    data.forEach((item: any) => {
      // Ensure item.geom exists and is valid
      if (item.geom) {
        const geoJsonLayer = L.geoJSON(item.geom, {
          title: item.TableName ,
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
    console.error("Error fetching GeoJSON data:", error);
    return []; // Return an empty array in case of error
  } finally {
    console.log("Map projection:", 'EPSG:3857'); // Leaflet uses Web Mercator by default
  }
}




async function getDataByTableClummnValue(path: string,table:string,column:string,operator:string,value:string, map: LeafletMap): Promise<any> {
  try {
    // Fetch GeoJSON data from the server
    const response = await fetch(`${serverbaseurl}/map/${path}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({selectedLayer:table,selectedAttribute:column,selectedOperator:operator,searchValue:value}) // Send selectedLayer in the body
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

    setResults(data.length)

    seTableData(data);
    setIsTableViewOpen(true)

    
    

    // Set table column names from the first record (you can adjust this logic)
    const columns = Object.keys(data[0]).map((key) => ({
      name: formatKey(key)
    }));

    setTableColumnNames(columns);


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
          title: item.TableName ,
          style: function (feature) {
            return { color: 'red'}; // You can customize the style here
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
    console.error("Error fetching GeoJSON data:", error);
    return []; // Return an empty array in case of error
  } finally {
    console.log("Map projection:", 'EPSG:3857'); // Leaflet uses Web Mercator by default
  }
}





async function getDataByTable(path: string,table:string, map: LeafletMap): Promise<any> {
  try {
    // Fetch GeoJSON data from the server
    const response = await fetch(`${serverbaseurl}/map/${path}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({selectedLayer:table}) // Send selectedLayer in the body
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

    setResults(data.length)

    seTableData(data);
    setIsTableViewOpen(true)

    
    

    // Set table column names from the first record (you can adjust this logic)
    const columns = Object.keys(data[0]).map((key) => ({
      name: formatKey(key)
    }));

    setTableColumnNames(columns);

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
        console.log("name",item.TableName)
        const geoJsonLayer = L.geoJSON(item.geom, {
          title: item.TableName ,
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
    console.error("Error fetching GeoJSON data:", error);
    return []; // Return an empty array in case of error
  } finally {
    console.log("Map projection:", 'EPSG:3857'); // Leaflet uses Web Mercator by default
  }
}




// Usage:
// getData('your_layer_name', mapInstance)
//   .then(features => {
//     if (features) {
//       console.log('Features loaded:', features.length);
//     }
//   })
//   .catch(error => console.error('Error:', error));

//  const filterQuery = useCallback((data:any) => {
//     if (!map) return;

//     if (geojsonLayer) {
//       map.removeLayer(geojsonLayer);
//     }

   

  

//     fetch(`${serverbaseurl}/map/featureoperatures`, {
//       method: 'POST',
//       credentials: 'include',
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify({ selectedAttribute:attribute,selectedLayer:selectedFeatureLayer}) // Send searchLayer in the body
//     })
//       .then(response => {
//         if (!response.ok) {
//           throw new Error(`HTTP error! status: ${response.status}`);
//         }
//         return response.json();
//       })
//       .then(data => {
//         console.log("GeoJSON data received:", data);

        

//         const geoJsonData = L.geoJSON(data, {
//           style: function (feature) {
//             return { color: 'blue' }; // You can customize the style here
//           },
//           onEachFeature: function (feature, layer) {
//             // You can add popups or tooltips here
//             if (feature.properties) {
//               layer.bindPopup(Object.keys(feature.properties).map(key => 
//                 feature.properties[key]!==null?`${key}: ${feature.properties[key]}`:null
//               ).join('<br>'));
//             }
//           }
//         });

//         console.log("Features loaded:", data.features.length);
//         setResults(data.features.length)

//         if (data.features.length > 0) {
//           // seTableData(data.features)

//           addLayerToMap(geoJsonData)
         
//           setGeojsonLayer(geoJsonData);

//           const extent = geoJsonData.getBounds();
//           const latLngBounds = [
//             [extent._southWest.lat, extent._southWest.lng], // Southwest corner
//             [extent._northEast.lat, extent._northEast.lng]  // Northeast corner
//         ]
//           console.log("extends",extent)
//           console.log("Layer extent:", latLngBounds);
//           if (extent) {
//             setBounds(latLngBounds)
//           } else {
//             console.log("Invalid or empty extent");
//           }

          
//         } 
//         else {
//           console.log("No features found for the query.");
//         }
//       })
//       .catch(error => {
//         console.error("Error fetching GeoJSON data:", error);
//       });

//     // console.log("Map projection:", map.getView().getProjection().getCode());
//   }, [map, geojsonLayer, selectedFeatureLayer, selectedAttribute, selectedQueryType, value]);



 







// Example usage








 


  const  handleQuery = async() => {
   
    if(selectedFeatureLayer ==="All Feature Layers" && selectedAttribute ==="None" && selectedQueryType ==="None"){
        console.log("Clicked")
        // tableData.length = 0
        // tableColumnNames.length = 0
      // const allAttributes= []
       
      // for (let i = 0; i < layers.length; i++) { const searchItem = "your_search_term"; // This could come from state or a prop
        // let newLayer = createWMSLayer(layers[i].name);
        // console.log("layer", newLayer);
        // addLayerToMap(newLayer);
        await getAllData("searchallfeaturelayersdata",map)
        // console.log("server data",data)
        
        // Uncomment and adjust these lines if you need to fetch additional data
        // var columnData = await getAttributeNames("nurc:"+layers[i].name);
        // allAttributes.push(...columnData);
        // await getData("nurc:"+layers[i].name);
      // }
        
        // setTableColumnNames(...new Set(allAttributes))
      
      //  await zoomToFitAllLayers() 
    }else if (selectedFeatureLayer !=="All Feature Layers" && selectedAttribute !=="None" && selectedQueryType ==="None"){

      await getDataByColumn("searchbycolumn",selectedFeatureLayer,selectedAttribute, map)
    
   
    }else if (selectedFeatureLayer !=="All Feature Layers" && selectedAttribute ==="None" && selectedQueryType ==="None"){
      await getDataByTable("searchbyfeaturelayer",selectedFeatureLayer, map)
      // tableData.length = 0
      
      
      // map.addLayer(await createWMSLayer(selectedFeatureLayer));
    //   var newLayer = createWMSLayer(selectedFeatureLayer)
    //   await getAttributeNames("nurc:"+selectedFeatureLayer)
      // setTableColumnNames(attributes)
    //   const data = await getData("nurc:"+selectedFeatureLayer)
      
    //   console.log("getdata",data)
      // setIsTableViewOpen(true)
    //   addLayerToMap(newLayer)
      // map.fitBounds(bounds)
      
      // layerSwitcherRef.current.addOverlay(newLayer)
      

        // Zoom to the extent of the new WMS layer
        // const extent = await  getWmsLayerExtent(selectedFeatureLayer);
        // console.log(extent)
       
        // const bounds = [
        //   [extent.southBoundLatitude, extent.westBoundLongitude], // Southwest corner
        //   [extent.northBoundLatitude, extent.eastBoundLongitude]  // Northeast corner
        // ];
        // map.fitBounds(bounds)
      
        // if (extent) {
        //     map.getView().fit(extent, {
        //         duration: 1000, // Animation duration
        //         Zoom: 20, // Optional: Set a maximum zoom level
        //     });
        //   }
      }else if (selectedFeatureLayer !=="All Feature Layers" && selectedAttribute !=="None" && selectedQueryType !=="None"){
        // tableData.length = 0
       
        await getDataByTableClummnValue("makeqquery",selectedFeatureLayer,selectedAttribute,selectedQueryType,value,map)
        // const data1 = await getAttributeNames("nurc:"+selectedFeatureLayer)
        // setTableColumnNames(data1)
        // const filterquery
        // setIsTableViewOpen(true)
        // filterQuery()
        
      }else if (selectedFeatureLayer ==="All Feature Layers" && selectedAttribute !=="None" && selectedQueryType ==="None"){
       console.log("Dont select selected attributes if selectedfeaturelayer is not none")
      }else if (selectedFeatureLayer ==="All Feature Layers" && selectedAttribute ==="None" && selectedQueryType !=="None"){
        console.log("Dont select selected selecedquerytype if selectedfeaturelayer is not none")
       }else if (selectedFeatureLayer !=="All Feature Layers" && selectedAttribute ==="None" && selectedQueryType ==="None"){
        console.log("Selecting None for both attribute and operators requires to select , All features Layers in the feature layer")
       }else {
        console.log(selectedFeatureLayer,selectedAttribute,selectedQueryType)
        console.log("unknown logic fix")
       }        
}

async function SetOperations(attribute:any){
  fetch(`${serverbaseurl}/map/featureoperatures`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ selectedAttribute:attribute,selectedLayer:selectedFeatureLayer}) // Send searchLayer in the body
  })
    .then(response_data => {
      if (!response_data.ok) {
        console.log(response_data.status);
        throw new Error('Network response was not ok');
      }
      return response_data.json();
    })
    .then(response_data => {
      // console.log("returned data", response_data.data);
      const newLayers = Object.values(response_data.data);
      setOperations(newLayers);
      console.log("operators",response_data.data );
    })
    .catch(error => {
      console.log('Unknown error occurred:', error.message);
    });
}

// useEffect(() => {
//   if (bounds) {
//       // Make sure bounds are valid before calling fitBounds
//       map.fitBounds(bounds);
//   }
// }, [bounds])

  // Log selected values whenever they change
  useEffect(() => {
    getAttributeNames(selectedFeatureLayer)
    // console.log("Selected Feature Layer:", selectedFeatureLayer);
    
  }, [selectedFeatureLayer]);

  useEffect(() => {
    
    SetOperations(selectedAttribute)
    
  }, [selectedAttribute]);

  // useEffect(() => {
  //   console.log("Selected Query Type:", selectedQueryType);
  // }, [selectedQueryType]);

  useEffect(() => {
   
  
    fetch(`${serverbaseurl}/map/featurelayers`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        "Content-Type": "application/json",
      },
    })
    .then(response_data => {
      if (!response_data.ok) {
        console.log(response_data.status);
        throw new Error('Network response was not ok');
      }
      return response_data.json();
    })
    .then(response_data => {
      console.log("returned data", response_data.data);
      const newLayers = Object.values(response_data.data);
      setLayers(newLayers);
      console.log("layers3", newLayers);
      // console.log(Object.values(response_data.data))
    })
    .catch(error => {
      console.error('Unknown error occurred:', error.message);
    });
  }, []); // Add searchItem to dependency array

  if (!isOpen) return null;


 

  return (
    <div className="flex w-1/4 flex-col mt-0 pt-0 h-3/4" >
      <div className="flex flex-col items-left bg-white p-2">
        <label htmlFor="feature-layer" className="mb-2">Select Feature Layer:</label>
        <select 
          id="feature-layer" 
          onChange={handleSelectedFeatureLayerOnChange} 
          value={selectedFeatureLayer} 
          className="border rounded p-2 mb-2"
        >
          <option value="All Feature Layers">All Feature Layers</option>
           {layers.map((layer, index) => (<option key ={index} value={layer}>{layer}</option>))}
          {/* <option value="All Feature Layers">All Feature Layers</option> */}
          {/* Add more options here */}
        </select>

        <label htmlFor="attribute" className="mb-2">Select Attribute:</label>
        <select 
          id="attribute" 
          onChange={handleSelectedAttributeOnChange} 
          value={selectedAttribute} 
          className="border rounded p-2 mb-2"
        >
          <option value="None">None</option>
          {attributes.map((attribute, index) => (<option key ={index} value={attribute}>{attribute}</option>))}
          {/* Add more options here */}
        </select>

        <label htmlFor="query-type" className="mb-2">Select Query Type:</label>
        <select 
          id="query-type" 
          onChange={handleSelectedQueryTypeOnChange} 
          value={selectedQueryType} 
          className="border rounded p-2 mb-2"
        >
          <option value="None">None</option>
          {operators.map((operator, index) => (<option key ={index} value={operator}>{operator}</option>))}
        
          {/* Add more options here */}
        </select>

        <label htmlFor="value" className="mb-2">Enter Value:</label>
        <input 
          id="value" 
          value={value}
          onChange={handleValueChange}
          className="border rounded p-2 mb-2"
        />

        <label htmlFor="results" className="mb-2">Number of Results Returned:</label>
        <input 
          id="results" 
          value={results}
          onChange={handleResultsChange}
          className="border rounded p-2 mb-2"
        />

<button 
onClick={handleQuery} 
className="mt-2" style={{ width: "60px", borderRadius: "2px", color: "white", background: "brown" }}>Query</button>
      </div>
    </div>
  );
}


export default QueryPanel
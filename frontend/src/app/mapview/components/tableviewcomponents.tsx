import React, { useEffect, useState, useRef } from 'react';
import L from 'leaflet';

const DetailsTableView = ({ tableData, tableColumnNames, map }) => {
  const [selectedFeatureId, setSelectedFeatureId] = useState(null);
  const tableRef = useRef(null);
  const highlightLayerRef = useRef(null);

  function formatValue(key: string, value: any): any {
    if (value && typeof value === 'object') {
      // Check if it's a GeoJSON-like object (or any object you expect)
      if (value.type && value.coordinates) {
        return 'GeoJSON Object';  // You can customize how you want to display it
      }
      // For general objects, you can either stringify or show some other representation
      return JSON.stringify(value);
    }
  
    const formatMap = {
      'shape__len': `${value} meters`,
      'shape__are': `${value} sq meters`,
      'creationda': value ? new Date(value).toLocaleDateString() : value
    };
  
    return formatMap[key] || value;
  }

  useEffect(() => {
    if (!map) return;

    highlightLayerRef.current = L.layerGroup().addTo(map);

    const handleMapClick = (event) => {
      map.eachLayer((mapLayer) => {
        if (mapLayer instanceof L.GeoJSON) {
          const clickedFeature = mapLayer.getLayers().find(feature => {
            if (feature instanceof L.Polygon || feature instanceof L.Polyline) {
              return feature.getBounds().contains(event.latlng);
            } else if (feature instanceof L.Marker) {
              return feature.getLatLng().equals(event.latlng);
            }
            return false;
          });

          if (clickedFeature) {
            highlightFeatureAndRow(clickedFeature);
          }
        }
      });
    };

    map.on('click', handleMapClick);

    map.eachLayer((mapLayer) => {
      if (mapLayer instanceof L.GeoJSON) {
        mapLayer.on('click', (e) => {
          L.DomEvent.stopPropagation(e);
          highlightFeatureAndRow(e.layer);
        });
      }
    });

    return () => {
      map.off('click', handleMapClick);
      if (highlightLayerRef.current) {
        map.removeLayer(highlightLayerRef.current);
      }
      map.eachLayer((mapLayer) => {
        if (mapLayer instanceof L.GeoJSON) {
          mapLayer.off('click');
        }
      });
    };
  }, [map]);

  const highlightFeatureAndRow = (feature) => {
    if (feature && highlightLayerRef.current) {
      highlightLayerRef.current.clearLayers();

      let highlightedFeature;
      if (feature instanceof L.Polygon || feature instanceof L.Polyline) {
        highlightedFeature = L.GeoJSON.geometryToLayer(feature.toGeoJSON());
      } else if (feature instanceof L.Marker) {
        highlightedFeature = L.marker(feature.getLatLng());
      }

      if (highlightedFeature) {
        highlightedFeature.setStyle({
          color: '#FF0000',
          weight: 3,
          opacity: 1,
          fillOpacity: 0.3
        });

        highlightLayerRef.current.addLayer(highlightedFeature);
      }

      if (feature.getBounds) {
        map.fitBounds(feature.getBounds(), { padding: [50, 50] });
      } else if (feature instanceof L.Marker) {
        map.setView(feature.getLatLng(), map.getZoom());
      }

      const featureId = feature.feature.properties.id || feature.feature.id;
      setSelectedFeatureId(featureId);
      scrollToTableRow(featureId);
    }
  };

  const scrollToTableRow = (featureId) => {
    if (tableRef.current) {
      const rowElement = tableRef.current.querySelector(`[data-feature-id="${featureId}"]`);
      if (rowElement) {
        rowElement.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
      }
    }
  };

  const handleRowClick = (feature) => {
    map.eachLayer((layer) => {
      if (layer instanceof L.GeoJSON) {
        const matchingFeature = layer.getLayers().find(f => 
          (f.feature.properties.id || f.feature.id) === (feature.properties.id || feature.id)
        );
        if (matchingFeature) {
          highlightFeatureAndRow(matchingFeature);
        }
      }
    });
  };

  return (
    <div className="h-1/4 z-[5000] bg-gray-800 w-full overflow-hidden relative">
  {/* Table Header */}
  <div className="overflow-hidden">
    <table className="text-white w-full border-collapse">
      <thead className="bg-gray-900 sticky top-0 z-10">
        <tr>
          {tableColumnNames
            .filter((column) => !column.name.toLowerCase().includes("geom")) // Filter out columns containing "Geom"
            .map((column, index) => (
              <th key={index} className="border border-white p-2">
                {column.name}
              </th>
            ))}
        </tr>
      </thead>
    </table>
  </div>

      {/* Scrollable Table Body */}
      <div className="overflow-y-auto h-full" ref={tableRef} style={{ maxHeight: 'calc(100% - 40px)' }}>
        <table className="text-white w-full border-collapse">
          <tbody>
            {tableData.map((feature, rowIndex) => {
              const featureId = feature?.properties?.id || feature?.id; // Check feature properties for id
              if (!featureId) return null; // Skip if no feature id

              return (
                <tr
                  key={rowIndex}
                  data-feature-id={featureId}
                  className={`
                    ${rowIndex % 2 === 0 ? 'bg-gray-100' : 'bg-gray-200'} 
                    hover:bg-gray-300
                    ${selectedFeatureId === featureId ? 'bg-yellow-200' : ''}
                  `}
                  onClick={() => handleRowClick(feature)}
                  style={{ cursor: 'pointer' }}
                >
                {Object.keys(feature)
                .filter((key) => !key.toLowerCase().includes("geom")) // Filter out data columns containing "Geom"
                .map((key, columnIndex) => (
                  <td key={`${featureId}-${columnIndex}`} className="border text-black border p-2">
                    {key !== undefined ? formatValue(key, feature[key]) : 'N/A'}
                  </td>
                ))}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DetailsTableView;

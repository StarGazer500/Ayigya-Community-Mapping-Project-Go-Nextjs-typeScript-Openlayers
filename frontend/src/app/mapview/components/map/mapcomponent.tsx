"use client";

import { MapContainer, TileLayer,  LayersControl } from "react-leaflet";
import { LatLngExpression, LatLngTuple } from 'leaflet';
import "leaflet/dist/leaflet.css";
import "leaflet-defaulticon-compatibility/dist/leaflet-defaulticon-compatibility.css";
import "leaflet-defaulticon-compatibility";
import { useEffect,useState, useRef } from "react";
import { Map as LeafletMap } from "leaflet"; 

import QueryPanel from "@/app/mapview/components/querypanelcomponent"
import SimpleQuery from "@/app/mapview/components/simplesearchcomponent"
import DetailsTableView  from "@/app/mapview/components/tableviewcomponents"
// import LayerSwitcher from "@/app/mapview/components/layerswitcher"

const { BaseLayer } = LayersControl;

// Define the interface for the Map component props
interface MapProps {
  posix: LatLngExpression | LatLngTuple;
  zoom?: number;
}

const defaults = {
  zoom: 19,
};

const Map = ({ posix, zoom = defaults.zoom }: MapProps) => {
  const mapRef = useRef<any>(null);
  const [marker, setMarker] = useState<any>(null);
  const [isOpen, setIsOpen] = useState<boolean>(false)

  const [tableColumnNames,setTableColumnNames] = useState<any[]>([])
  const [tableData,seTableData] = useState<any[]>([])
  const [results, setResults] = useState<string>('');
  const [isTableViewOpen,setIsTableViewOpen] = useState<boolean>(false)
  const [simpleQueryResults,setSimpleQueryResults] = useState<any>(null)

  useEffect(()=>{
    if (mapRef.current && mapRef.current instanceof LeafletMap){

     console.log("map initiated")    
    }
  },[mapRef.current])

  const addLayerToMap = (layer:any):void => {
    if (mapRef.current && layer && !mapRef.current.hasLayer(layer)) {
      mapRef.current.addLayer(layer);
      // mapRef.current.invalidateSize();
      console.log(`Layer ${layer.options.title} added to map`);
    }
  };

  
  return (
    <div className="flex flex-col w-full h-screen" >

      <div style={{ background: 'brown' }} className="mb-0 pb-0">
          <h1 className="font-bold text-6xl pb-0 mb-0 text-white">Ayigya Community Map View</h1>
      </div>
<     div className="flex flex-col flex-1 overflow-hidden">
<       div className="flex flex-1">
      <QueryPanel addLayerToMap = {addLayerToMap}  setIsTableViewOpen ={setIsTableViewOpen} results = {results} setResults = {setResults} tableData = {tableData} tableColumnNames={tableColumnNames} setTableColumnNames = {setTableColumnNames} seTableData = {seTableData}  isOpen={isOpen} map={mapRef.current} />
      
    <div className="flex-1 relative">
    {/* <LayerSwitcher map = {mapRef.current}/> */}
    <button
              style={{ zIndex: 4000 }}
              onClick={() => setIsOpen(!isOpen)}
              className="absolute top-2 left-20 p-1 rounded-md text-md text-white font-bold bg-[#CD4631]"
            >
              {isOpen ? "Hide Advanced Query" : "Open Advanced Query"}
            </button>
            <SimpleQuery  addLayerToMap = {addLayerToMap}  setIsTableViewOpen = {setIsTableViewOpen} setSimpleQueryResults = {setSimpleQueryResults} setTableColumnNames = {setTableColumnNames} seTableData = {seTableData} map={mapRef.current} />
            <div className="h-full bg-grey relative">
            <MapContainer 
              center={posix} 
              zoom={20} 
              // style={{ height: "100vh", width: "100%" }} 
              className="w-full h-full  bg-grey"
              scrollWheelZoom={true}
              ref = {mapRef}
            
             
            >
              <LayersControl position="topright">
                <BaseLayer  name="OpenStreetMap">
                  <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                  />
                </BaseLayer>
                <BaseLayer name="Google Satellite">
                  <TileLayer
                    url="https://mt1.google.com/vt/lyrs=s&x={x}&y={y}&z={z}"
                    attribution='&copy; Google'
                  />
                </BaseLayer>
                <BaseLayer checked name="Google Satellite Hybrid">
                  <TileLayer
                    url="https://mt1.google.com/vt/lyrs=y&x={x}&y={y}&z={z}"
                    attribution='&copy; Google'
                  />
                </BaseLayer>
              </LayersControl>
            </MapContainer>
            </div>
             {/* <Legend  map={mapRef.current} />  */}
            {simpleQueryResults?<div ><h2 style={{color:"white"}} className='absolute z-[9000] font-bold  bg-red-500 p-1 rounded-md  bottom-6 left-2'>{simpleQueryResults} Results Returned</h2></div>:null} 
            </div>
        </div>
        {isTableViewOpen?<DetailsTableView map = {mapRef.current} tableColumnNames = {tableColumnNames} tableData= {tableData}/>:null}
      </div>
    </div>
  );
};

export default Map;

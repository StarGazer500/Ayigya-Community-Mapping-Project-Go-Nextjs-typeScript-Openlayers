"use client";

import dynamic from "next/dynamic";
import { Suspense } from "react";

// Dynamically import the Map component with no SSR
const Map = dynamic(
  () => import('@/app/mapview/components/map/mapcomponent'),
  {
    loading: () => <p>A map is loading...</p>,
    ssr: false, // Disable SSR for this dynamic import
  }
);

export default function Page() {

  
  return (
    <div className="bg-white-700 mx-auto my-5 w-[98%] h-[480px]">
      <Suspense fallback={<p>Loading map...</p>}>
        {/* Pass latitude and longitude coordinates to the Map component */}
        <Map posix={[4.79029, -75.69003]} />
      </Suspense>
    </div>
  );
}

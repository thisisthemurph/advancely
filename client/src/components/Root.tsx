import { Outlet } from "react-router-dom";
import Header from "./Header";

function Root() {
  return (
    <>
      <Header />
      <main className="p-8">
        <Outlet />
      </main>
    </>
  );
}

export default Root;

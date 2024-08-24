import { createBrowserRouter } from "react-router-dom"
import Root from "../components/Root.tsx";
import HomePage from "../pages/home/index.tsx";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    children: [
      {
        path: "/",
        element: <HomePage />
      }
    ]
  }
]);
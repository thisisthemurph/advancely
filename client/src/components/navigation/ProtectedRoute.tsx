import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth.tsx";

const ProtectedRoute = () => {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }
  return <Outlet />
}

export default ProtectedRoute;

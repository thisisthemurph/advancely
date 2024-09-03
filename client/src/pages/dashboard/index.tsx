import PageHeading from "../../components/ui/PageHeading";
import { useAuth } from "../../hooks/useAuth.tsx";

function DashboardPage() {
  const auth = useAuth();

  return (
    <>
      <PageHeading heading="Dashboard" />
      <p>Welcome to your dashboard</p>
      <p>{auth.isAuthenticated ? "Authenticated" : "NOT AUTHENTICATED"}</p>
      <p>{auth.session?.email}</p>
    </>
  );
}

DashboardPage.path = "/dashboard";

export default DashboardPage;

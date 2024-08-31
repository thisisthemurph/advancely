import PageHeading from "../../components/ui/PageHeading";

function DashboardPage() {
  return (
    <>
      <PageHeading heading="Dashboard" />
      <p>Welcome to your dashboard</p>
    </>
  );
}

DashboardPage.path = "/dashboard";

export default DashboardPage;

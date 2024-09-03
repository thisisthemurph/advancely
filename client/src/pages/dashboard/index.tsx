import PageHeading from "../../components/ui/PageHeading";
import { useAuth } from "../../hooks/useAuth.tsx";

function DashboardPage() {
  const auth = useAuth();

  const getGreeting = (): string => {
    const hour = new Date().getHours();

    if (hour < 12) {
      return "Good morning";
    } else if (hour < 18) {
      return "Good afternoon";
    }
    return "Good evening";
  }

  return (
    <>
      <section className="mb-8 space-y-2">
        <PageHeading heading="Dashboard" />
        <p className="text-lg font-semibold">{getGreeting()}, {auth.user?.firstName}!</p>
      </section>
    </>
  );
}

DashboardPage.path = "/dashboard";

export default DashboardPage;

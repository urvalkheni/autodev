import Navbar from "@/components/Navbar";
import Hero from "@/components/Hero";
import Terminal from "@/components/Terminal";
import Categories from "@/components/Categories";
import Profiles from "@/components/Profiles";
import GithubScanner from "@/components/GithubScanner";
import Skills from "@/components/Skills";
import InstallMethods from "@/components/InstallMethods";
import Footer from "@/components/Footer";

export default function Home() {
  return (
    <main className="min-h-screen bg-black">
      <Navbar />
      <Hero />
      <Terminal />
      <Categories />
      <Profiles />
      <GithubScanner />
      <Skills />
      <InstallMethods />
      <Footer />
    </main>
  );
}

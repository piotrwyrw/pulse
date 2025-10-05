import {bootstrapApplication} from '@angular/platform-browser';
import {appConfig} from './app/app.config';
import {AppComponent} from './app/app.component';
import {
  BarController,
  BarElement,
  CategoryScale,
  Chart, Legend,
  LinearScale,
  LineController,
  LineElement,
  PointElement, Title, Tooltip
} from "chart.js";

Chart.register(
  BarController,
  LineController,
  PointElement,
  LineElement,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);
Chart.defaults.backgroundColor = "#1b1b1b";
Chart.defaults.borderColor = "#ffffff";
Chart.defaults.color = "#ffffff";

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));

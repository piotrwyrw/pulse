import { Component } from '@angular/core';
import {PlotComponent} from "../../components/plot/plot.component";

@Component({
  selector: 'app-home',
  imports: [
    PlotComponent
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss'
})
export class HomeComponent {

}

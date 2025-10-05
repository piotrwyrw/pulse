import {Component, ElementRef, Input, ViewChild} from '@angular/core';
import {Chart} from "chart.js";
import {PulseService} from "../../services/pulse.service";

@Component({
  selector: 'app-plot',
  imports: [],
  templateUrl: './plot.component.html',
  styleUrl: './plot.component.scss'
})
export class PlotComponent {
  @ViewChild('plot_canvas') plot!: ElementRef<HTMLCanvasElement>;

  @Input() color: string = 'rgb(162, 155, 254)';

  constructor(public pulseService: PulseService) {
  }

  private ngAfterViewInit() {
    new Chart(this.plot.nativeElement, {
      type: 'line',
      data: {
        labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
        datasets: [{
          label: '# of Votes',
          data: [12, 19, 3, 5, 2, 3],
          fill: true,
          borderWidth: 3,
          borderColor: this.color,
          cubicInterpolationMode: "monotone"
        }]
      },
      options: {
        scales: {
          x: {
            grid: {
              display: false
            }
          },
          y: {
            grid: {
              display: false
            }
          }
        },

      }
    });
  }


}

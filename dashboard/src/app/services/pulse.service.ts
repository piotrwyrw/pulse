import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../env/environment";
import {PulseRecord} from "../models/record";

@Injectable({
  providedIn: 'root'
})
export class PulseService {
  private data: PulseRecord[] = []

  constructor(private http: HttpClient) {
    this.fetchData()
  }

  fetchData() {
    this.http.get("http://" + environment.pulseServer + "/records").subscribe({
      next: data => {
        console.log("OK")
        console.log(data)
      },
      error: err => {
        console.log(err)
      }
    })
  }

}

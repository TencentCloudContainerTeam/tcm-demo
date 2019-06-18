package com.webdemo.controller;

import com.google.common.collect.Maps;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.math.BigDecimal;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

@RestController
public class SalesController {

    @GetMapping("/sales")
    public Map<Integer, Integer> sales(@RequestParam("ids") String ids) {
        Map<Integer, Integer> mockDbMap = Maps.newHashMap();
        mockDbMap.put(1, 502);
        mockDbMap.put(2, 408);
        mockDbMap.put(3, 993);
        mockDbMap.put(4, 981);
        mockDbMap.put(5, 270);
        mockDbMap.put(6, 345);
        mockDbMap.put(7, 89);
        mockDbMap.put(8, 44);
        mockDbMap.put(9, 290);
        mockDbMap.put(10, 88);
        mockDbMap.put(11, 60);
        mockDbMap.put(12, 51);
        mockDbMap.put(13, 302);
        mockDbMap.put(14, 49);
        mockDbMap.put(15, 91);

        Map<Integer, Integer> resultMap = Maps.newHashMap();
        System.out.println("getting sales of ids: " + ids);

        for (String id : ids.split(",")) {
            Integer value = mockDbMap.get(Integer.valueOf(id));
            if (value != null) {
                // BigDecimal multiply = new BigDecimal(value).multiply(new BigDecimal(10));
                // resultMap.put(Integer.valueOf(id), multiply.toString());
                resultMap.put(Integer.valueOf(id), value);
            }
        }

        return resultMap;
    }

}

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
public class StockController {

    @GetMapping("/stock")
    public Map<Integer, Integer> sales(@RequestParam("ids") String ids) {
        Map<Integer, Integer> mockDbMap = Maps.newHashMap();
        mockDbMap.put(1, 2001);
        mockDbMap.put(2, 3002);
        mockDbMap.put(3, 1903);
        mockDbMap.put(4, 1404);
        mockDbMap.put(5, 2005);
        mockDbMap.put(6, 1000);
        mockDbMap.put(7, 1007);
        mockDbMap.put(8, 1008);
        mockDbMap.put(9, 1009);
        mockDbMap.put(10, 1010);
        mockDbMap.put(11, 1511);
        mockDbMap.put(12, 1612);
        mockDbMap.put(13, 1713);
        mockDbMap.put(14, 1814);
        mockDbMap.put(15, 1015);

        Map<Integer, Integer> resultMap = Maps.newHashMap();

        System.out.println("getting stock of ids: " + ids);
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
